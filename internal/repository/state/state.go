package state

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	stateentity "github.com/albertwidi/go-project-example/internal/entity/state"
	"github.com/albertwidi/go-project-example/internal/pkg/redis"
	"github.com/albertwidi/go-project-example/internal/pkg/ulid"
)

// Repository struct
type Repository struct {
	Redis redis.Redis
	// ulid generator
	ulidgen *ulid.Ulid
}

// New repository
func New(redis redis.Redis) *Repository {
	r := Repository{
		Redis:   redis,
		ulidgen: ulid.New(3),
	}

	return &r
}

func formatStateIDKey(id string) string {
	return fmt.Sprintf("state:%s", id)
}

// Save state
func (r *Repository) Save(ctx context.Context, id string, data stateentity.State) error {
	stateID := formatStateIDKey(id)

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = r.Redis.SetEX(ctx, stateID, string(out), int(data.ExpiryTime.Seconds()))
	if err != nil {
		return err
	}

	return nil
}

// Get return available state
func (r *Repository) Get(ctx context.Context, stateID string) (stateentity.State, error) {
	stt := stateentity.State{}

	stateID = formatStateIDKey(stateID)
	val, err := r.Redis.Get(stateID)
	if err != nil {
		if redis.IsErrNil(err) {
			return stt, stateentity.ErrStateNotFound
		}

		return stt, err
	}

	if val == "" {
		return stt, stateentity.ErrStateNotFound
	}

	err = json.Unmarshal([]byte(val), &stt)
	if err != nil {
		return stt, err
	}

	return stt, nil
}

// SetExpire to set expire for state
func (r *Repository) SetExpire(ctx context.Context, stateID string, expiryTime time.Duration) error {
	// op := xerrors.EOp("stateResource/setExpire")

	stateID = formatStateIDKey(stateID)
	_, err := r.Redis.Expire(stateID, int(expiryTime.Seconds()))
	if err != nil {
		return err
	}

	return err
}

// Delete state
func (r *Repository) Delete(ctx context.Context, stateID string) error {
	stateID = formatStateIDKey(stateID)
	_, err := r.Redis.Delete(stateID)
	if err != nil {
		return err
	}

	return err
}
