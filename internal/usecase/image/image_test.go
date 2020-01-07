package image_test

import (
	"context"
	"errors"
	"io"
	"testing"

	imageentity "github.com/albertwidi/go-project-example/internal/entity/image"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage/local"
	imageusecase "github.com/albertwidi/go-project-example/internal/usecase/image"
	imagemock "github.com/albertwidi/go-project-example/internal/usecase/image/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func newLocalStorage(bucketName string) (*objectstorage.Storage, error) {
	storage, err := local.New(context.Background(), bucketName, &local.Options{DeleteOnClose: true})
	if err != nil {
		return nil, err
	}
	return objectstorage.New(storage), nil
}

func TestUpload(t *testing.T) {
	t.Parallel()

	repo := imagemock.NewMockimageRepository(gomock.NewController(t))
	repo.EXPECT().
		GetTempPath(context.Background(), "abcd").
		Return("jkl", nil)
	repo.EXPECT().
		GetTempPath(context.Background(), "jkl").
		Return("", errors.New("fuck you"))

	storage, err := newLocalStorage("./testUpload")
	if err != nil {
		t.Error(err)
		return
	}
	defer storage.Close()

	usecase, err := imageusecase.New(storage, repo, nil)
	if err != nil {
		t.Error(err)
		return
	}

	cases := []struct {
		reader      io.Reader
		info        imageentity.FileInfo
		expectImage imageusecase.Image
		expectError error
	}{
		{
			info: imageentity.FileInfo{
				FileName: "testing.go",
				Size:     1000,
				UserHash: "eUjks",
				Mode:     imageentity.ModePrivate,
				Group:    imageentity.GroupUserAvatar,
				Tags:     "asd,jkl,abcd",
			},
		},
	}

	for _, c := range cases {
		img, err := usecase.Upload(context.Background(), c.reader, c.info)
		if err != c.expectError {
			t.Errorf("testUpload: expecting error %v but got %v", c.expectError, err)
			return
		}

		if !cmp.Equal(img, c.expectImage) {
			t.Errorf("testUpload: expecting %+v but got %+v", c.expectImage, img)
			return
		}
	}
}

func TestDownload(t *testing.T) {
}
