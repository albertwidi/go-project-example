package state

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	entity "github.com/albertwidi/go-project-example/internal/entity/state"
	"github.com/albertwidi/go-project-example/internal/pkg/ulid"
)

// Usecase of state
type Usecase struct {
	repo stateRepository
	// ulid generator
	ulidgen *ulid.Ulid
}

type stateRepository interface {
	Save(ctx context.Context, stateid string, data entity.State) error
	Get(ctx context.Context, stateid string) (entity.State, error)
	Delete(ctx context.Context, stateid string) error
	SetExpire(ctx context.Context, stateid string, duration time.Duration) error
}

// New state Usecase
func New(stateRepo stateRepository) *Usecase {
	u := Usecase{
		repo:    stateRepo,
		ulidgen: ulid.New(3),
	}
	return &u
}

// create id for state
func (u Usecase) createID(data entity.State) (string, error) {
	now := time.Now()

	// key for state_id generation conains ulid:identifier:timeunix
	keyList := []string{u.ulidgen.Ulid(), data.Identifier, now.String()}
	key := strings.Join(keyList, ":")
	// encode string to create a stateID
	stateID := base64.RawStdEncoding.EncodeToString([]byte(key))

	return stateID, nil
}

// Create state
// returning state_id and error
func (u Usecase) Create(ctx context.Context, data entity.State) (string, error) {
	if err := data.Validate(); err != nil {
		return "", err
	}

	id, err := u.createID(data)
	if err != nil {
		return "", err
	}

	if id == "" {
		return "", errors.New("state id is empty")
	}

	now := time.Now()
	if data.ExpiryTime == 0 {
		data.ExpiryTime = entity.DefaultStateExpiryTime
	}

	// set the expire at of the data
	data.ExpiredAt = now.Add(data.ExpiryTime)
	if err := u.repo.Save(ctx, id, data); err != nil {
		return "", err
	}

	return id, nil
}

// Get state
func (u Usecase) Get(ctx context.Context, stateid string) (entity.State, error) {
	state, err := u.repo.Get(ctx, stateid)
	return state, err
}

// Delete state
func (u Usecase) Delete(ctx context.Context, stateid string) error {
	err := u.repo.Delete(ctx, stateid)
	return err
}

// SetExpire for state
func (u Usecase) SetExpire(ctx context.Context, stateid string, duration time.Duration) error {
	return nil
}
