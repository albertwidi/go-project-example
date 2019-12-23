package user

import (
	"context"

	userentity "github.com/albertwidi/go-project-example/internal/entity/user"
	"github.com/albertwidi/go-project-example/internal/pkg/redis"
	"github.com/albertwidi/go-project-example/internal/pkg/sqldb"
)

// Repository of user
type Repository struct {
	db    *sqldb.DB
	redis redis.Redis
}

// New repository of user
func New(db *sqldb.DB, redis redis.Redis) *Repository {
	r := Repository{
		db:    db,
		redis: redis,
	}
	return &r
}

// Create user
func (r *Repository) Create(ctx context.Context, user userentity.User) (string, error) {
	return "", nil
}

// Update user
func (r *Repository) Update(ctx context.Context) error {
	return nil
}
