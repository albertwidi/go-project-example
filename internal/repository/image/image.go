package image

import (
	"context"
	"strings"
	"time"

	imageentity "github.com/albertwidi/go-project-example/internal/entity/image"
	"github.com/albertwidi/go-project-example/internal/pkg/redis"
)

// Repository of image
type Repository struct {
	redis redis.Redis
}

// New repository for image
func New(redis redis.Redis) *Repository {
	r := Repository{
		redis: redis,
	}
	return &r
}

func createImageKey(id string) string {
	return strings.Join([]string{"image_temp", id}, ":")
}

// SaveTempPath image path
func (r Repository) SaveTempPath(ctx context.Context, id, originalPath string, expiryTime time.Duration) error {
	key := createImageKey(id)
	_, err := r.redis.SetEX(ctx, key, originalPath, int(expiryTime.Seconds()))
	return err
}

// GetTempPath will return the original path from a temporary id
func (r Repository) GetTempPath(ctx context.Context, id string) (string, error) {
	key := createImageKey(id)
	out, err := r.redis.Get(ctx, key)
	if err != nil {
		if r.redis.IsErrNil(err) {
			return "", imageentity.ErrTempPathNotFound
		}
		return "", err
	}
	return out, err
}
