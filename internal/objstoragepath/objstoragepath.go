package objstoragepath

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"
	"time"

	imageentity "github.com/albertwidi/go-project-example/internal/entity/image"
)

// Config of image path library
type Config struct {
	Public  DownloadConfig
	Private DownloadConfig
}

// DownloadConfig struct
type DownloadConfig struct {
	DownloadProto string
	DownloadHost  string
	DownloadPort  string
	DownloadPath  string
}

// FilePath struct
type FilePath struct {
	Proto        string
	Host         string
	FilePath     string
	DownloadPath string
	DownloadLink string
	Signed       bool
}

// ObjectStoragePath generator
type ObjectStoragePath struct {
	config *Config
}

// New object storage path
func New(c *Config, local bool) (*ObjectStoragePath, error) {
	if c == nil {
		return nil, errors.New("imagepath: config cannot be nil")
	}

	if local {
		dialer := net.Dialer{Timeout: time.Second * 3}
		// TODO: check when internet connection is down
		conn, err := dialer.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return nil, err
		}

		addr := conn.LocalAddr().(*net.UDPAddr)
		// split the port
		s := strings.Split(addr.String(), ":")
		c.Public.DownloadHost = s[0]
		c.Private.DownloadHost = s[0]
	}
	return &ObjectStoragePath{config: c}, nil
}

// GetDownloadPath return needed download path
func (o *ObjectStoragePath) GetDownloadPath(mode imageentity.Mode) (string, error) {
	switch mode {
	case imageentity.ModePublic:
		return o.config.Public.DownloadPath, nil
	case imageentity.ModePrivate:
		return o.config.Private.DownloadPath, nil
	default:
		return "", fmt.Errorf("objectstoragepath: invalid mode, got %s", mode)
	}
}

// Generate return file path struct with all download information needed
func (o *ObjectStoragePath) Generate(mode imageentity.Mode, filePath string) (*FilePath, error) {
	var file FilePath
	if err := mode.Validate(); err != nil {
		return nil, err
	}

	switch mode {
	case imageentity.ModePublic:
		u, err := url.Parse(o.config.Public.DownloadProto + o.config.Public.DownloadHost + o.config.Public.DownloadPort)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, o.config.Public.DownloadPath, filePath)

		file = FilePath{
			Proto:        o.config.Public.DownloadProto,
			Host:         o.config.Public.DownloadHost,
			FilePath:     filePath,
			DownloadPath: path.Join(o.config.Public.DownloadPath, filePath),
			DownloadLink: u.String(),
		}

	case imageentity.ModePrivate:
		u, err := url.Parse(o.config.Private.DownloadProto + o.config.Private.DownloadHost + o.config.Private.DownloadPort)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, o.config.Private.DownloadPath)

		v := u.Query()
		v.Add("image_path", filePath)
		u.RawQuery = v.Encode()

		file = FilePath{
			Proto:        o.config.Private.DownloadProto,
			Host:         o.config.Private.DownloadHost,
			FilePath:     filePath,
			DownloadPath: path.Join(o.config.Private.DownloadPath, filePath),
			DownloadLink: u.String(),
		}

	case imageentity.ModeSigned:
		file = FilePath{
			Proto:        "",
			Host:         "",
			FilePath:     "",
			DownloadPath: "",
			DownloadLink: filePath,
			Signed:       true,
		}
	}
	return &file, nil
}
