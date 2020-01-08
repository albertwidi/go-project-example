package image_test

import (
	"context"
	"testing"
	"time"

	"github.com/albertwidi/go-project-example/internal/repository/image"
	redismock "github.com/albertwidi/go-project-example/internal/pkg/redis/mock"
	"github.com/golang/mock/gomock"
)

func TestSaveTempPath(t *testing.T) {
	t.Parallel()

	redisMock := redismock.NewMockRedis(gomock.NewController(t))
	redisMock.EXPECT().
		SetEX(context.Background(), "image_temp:abcd", "jklf", int(time.Minute.Seconds())).
		Return("OK", nil)

	repo := image.New(redisMock)
	cases := []struct {
		key          string
		value        string
		expiryTime   time.Duration
		expectResult string
		expectError  error
	}{
		{
			key:         "abcd",
			value:       "jklf",
			expiryTime:  time.Minute,
			expectError: nil,
		},
	}

	for _, c := range cases {
		err := repo.SaveTempPath(context.Background(), c.key, c.value, c.expiryTime)
		if err != c.expectError {
			t.Errorf("saveTempPath: expect error %v but got %v", c.expectError, err)
			return
		}
	}
}

func TestGetTempPath(t *testing.T) {

}