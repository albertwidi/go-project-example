package image

//go:generate mockgen -source=image.go -destination=mock/image_mock.go -package=image

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	imageentity "github.com/albertwidi/go-project-example/internal/entity/image"
	userentity "github.com/albertwidi/go-project-example/internal/entity/user"
	"github.com/albertwidi/go-project-example/internal/objstoragepath"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage"
	"github.com/albertwidi/go-project-example/internal/xerrors"
	guuid "github.com/google/uuid"
)

// Usecase of image
type Usecase struct {
	publicStorage  *objectstorage.Storage
	privateStorage *objectstorage.Storage
	imageRepo      imageRepository
	objPath        *objstoragepath.ObjectStoragePath
}

type imageRepository interface {
	SaveTempPath(ctx context.Context, id, originalPath string, expiryTime time.Duration) error
	GetTempPath(ctx context.Context, id string) (string, error)
}

// New image usecase
func New(privateStorage *objectstorage.Storage, imageRepo imageRepository, objStoragePath *objstoragepath.ObjectStoragePath) (*Usecase, error) {
	u := Usecase{
		privateStorage: privateStorage,
		imageRepo:      imageRepo,
	}
	return &u, nil
}

// createID for temporary path
func (u Usecase) createID(ctx context.Context, filepath string) (string, error) {
	// key is tempImagePath:{filepath}:{uuid}:{time_string}
	keys := []string{"tempImagePath", filepath, guuid.New().String(), time.Now().String()}
	key := strings.Join(keys, ":")

	// create id based on sha1
	s := sha1.New()
	s.Write([]byte(key))
	encryptedKey := hex.EncodeToString(s.Sum(nil))

	return encryptedKey, nil
}

// Image struct
type Image struct {
	Proto        string
	Host         string
	FilePath     string
	DownloadPath string
	DownloadLink string
}

// Upload image
func (u *Usecase) Upload(ctx context.Context, reader io.Reader, info imageentity.FileInfo) (Image, error) {
	image := Image{}
	if len(strings.Split(info.Tags, ",")) > 5 {
		return image, imageentity.ErrTooManyTags
	}
	if err := info.Mode.Validate(); err != nil {
		return image, err
	}

	var (
		filePath string
		err      error
		// file naming
		rename   bool
		fileName string
	)

	fileName = info.FileName
	group := info.Group
	if group == imageentity.GroupEmpty {
		group = imageentity.GroupMixed
	}

	if err := group.Validate(); err != nil {
		return image, err
	}

	// sometimes image needs rename to avoid naming collission in the bucket
	// TODO: it might be better to append the image name with timestamp on the server or special_id
	if rename {
		fileName = path.Join(info.FileName, time.Now().String())
		fileName = base64.RawStdEncoding.EncodeToString([]byte(fileName))
	}

	if info.Mode == imageentity.ModePrivate {
		key := path.Join(string(group), info.FileName)
		writeOptions := &objectstorage.WriteOptions{
			Metadata: map[string]string{
				"access_owner": string(imageentity.CreateAccess([]string{string(info.UserHash)}, []string{"read"})),
				"access_admin": string(imageentity.CreateAccess([]string{"all"}, []string{"read", "write"})),
				"tags":         info.Tags,
			},
		}

		filePath, err = u.privateStorage.Upload(ctx, reader, key, writeOptions)
		if err != nil {
			return image, err
		}

		// always generate temporary image path if user hash is empty
		if info.UserHash == "" {
			filePath, err = u.GenerateTemporaryPath(ctx, filePath, time.Minute*30)
			if err != nil {
				return image, err
			}
		}
	}

	img, err := u.objPath.Generate(info.Mode, filePath)
	if err != nil {
		return image, err
	}
	image = Image{
		Proto:        img.Proto,
		Host:         img.Host,
		FilePath:     img.FilePath,
		DownloadPath: img.DownloadPath,
		DownloadLink: img.DownloadLink,
	}
	return image, nil
}

// Download image
// download usecase never use a public storage
// as public storage being served via CDN and local filesystem in LOCAL
// to download temporary filepath, use format: temporary:{id}
func (u *Usecase) Download(ctx context.Context, imagePath string, userHash userentity.Hash, admin bool) ([]byte, error) {
	var (
		err error
	)

	prefix, filepath, err := u.GetImageFilePath(ctx, imagePath)
	if err != nil {
		return nil, err
	}

	// if the image path is not temporary, then we need to check the access previledge
	if prefix != prefixTemporary {
		attr, err := u.privateStorage.Attributes(ctx, filepath)
		if err != nil {
			return nil, err
		}

		var (
			access  []string
			granted bool
		)

		// choose one, so we don't have to do things twice
		if admin {
			accessData := attr.Metadata["access_admin"]
			access = strings.Split(accessData, ";")
		} else {
			accessData := attr.Metadata["access_owner"]
			access = strings.Split(accessData, ";")
		}

		for _, acc := range access {
			detail := strings.Split(acc, ":")
			if len(detail) < 2 {
				return nil, xerrors.New(imageentity.ErrInvalidAccessAttribute, xerrors.KindInternalError)
			}

			switch detail[0] {
			case "allowed":
				// allow immediately if admin
				if admin && detail[1] == "all" {
					granted = true
					break
				}

				if detail[1] == string(userHash) {
					granted = true
					break
				}

			case "priviledge":
				// not yet implemented
				// priviledge restrict access of read/write or deletion of the object
			}

			if granted {
				break
			}
		}
		if !granted {
			return nil, errors.New("image: user is not permitted to view this image/file")
		}
	}

	out, err := u.privateStorage.DownloadByte(ctx, filepath, nil)
	return out, err
}

// GenerateSignedURL for generating temporary path to download image directly from object storage provider
// this method is different from temporary as we are not serving the download from our server
func (u *Usecase) GenerateSignedURL(ctx context.Context, filePath string, expiry time.Duration) (string, error) {
	var (
		url string
		err error
	)

	// use temporary url if the storage is local
	if u.privateStorage.Name() == objectstorage.StorageLocal {
		url, err = u.GenerateTemporaryURL(ctx, filePath, expiry)
	} else {
		url, err = u.privateStorage.SignedURL(ctx, filePath, expiry)
	}

	if err != nil {
		return "", err
	}

	return url, nil
}

// GenerateTemporaryURL directly from backend
// instad of using object-storage url, this will use image-proxy url
func (u *Usecase) GenerateTemporaryURL(ctx context.Context, filePath string, expiry time.Duration) (string, error) {
	path, err := u.GenerateTemporaryPath(ctx, filePath, expiry)
	if err != nil {
		return "", err
	}

	img, err := u.objPath.Generate(imageentity.ModePrivate, path)
	if err != nil {
		return "", err
	}
	return img.DownloadLink, nil
}

// GenerateTemporaryPath for generating path to download image
// but with disposable image path
// this is useful for secure image handling
func (u *Usecase) GenerateTemporaryPath(ctx context.Context, filepath string, expiryTime time.Duration) (string, error) {
	id, err := u.createID(ctx, filepath)
	if err != nil {
		return "", err
	}

	if err := u.imageRepo.SaveTempPath(ctx, id, filepath, expiryTime); err != nil {
		return "", err
	}

	downloadPath, err := u.objPath.GetDownloadPath(imageentity.ModePrivate)
	if err != nil {
		return "", err
	}
	id = path.Join(downloadPath, id)
	return strings.Join([]string{"temporary", id}, ":"), nil
}

const (
	prefixTemporary = "temporary"
	prefixFile      = "file"
)

// GetImageFilePath normalize and return the true image file path needed by backend
func (u *Usecase) GetImageFilePath(ctx context.Context, imagePath string) (prefix, filePath string, err error) {
	s := strings.Split(imagePath, ":")
	slen := len(s)
	if slen < 1 {
		err = fmt.Errorf("imagepath: image path format not valid, got %s", imagePath)
		return
	}

	if slen == 1 {
		filePath = s[0]
		return
	}

	switch s[0] {
	case prefixTemporary:
		prefix = prefixTemporary
		filePath, err = u.filepathFromTemporaryPath(ctx, s[1])
		return
	case prefixFile:
		prefix = prefixFile
		filePath = s[1]
		return
	default:
		err = fmt.Errorf("imagepath: prefix is not valid, got %s", s[0])
		return
	}
}

// FilePathFromTemporaryPath path return the true filepath from temporary endpoint
func (u *Usecase) FilePathFromTemporaryPath(ctx context.Context, tempPath string) (string, error) {
	downloadPath, err := u.objPath.GetDownloadPath(imageentity.ModePrivate)
	if err != nil {
		return "", err
	}
	s := strings.SplitAfter(tempPath, downloadPath)
	if len(s) < 2 {
		return "", errors.New("temporary path is not valid")
	}

	filePath := s[1]
	filePath = strings.TrimLeft(filePath, "/")

	return u.filepathFromTemporaryPath(ctx, filePath)
}

// filepathFromTemporaryPath return the original path of temporary path
func (u *Usecase) filepathFromTemporaryPath(ctx context.Context, id string) (string, error) {
	path, err := u.imageRepo.GetTempPath(ctx, id)
	if err != nil {
		return "", err
	}

	return path, err
}
