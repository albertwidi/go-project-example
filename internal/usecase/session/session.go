package session

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	sessionentity "github.com/albertwidi/go-project-example/internal/entity/session"
	userentity "github.com/albertwidi/go-project-example/internal/entity/user"
	guuid "github.com/google/uuid"
)

// UseCase struct
type UseCase struct {
	repo sessionRepo
}

// sessionRepo interface
type sessionRepo interface {
	Save(ctx context.Context, userhash userentity.Hash, sessionid string, data sessionentity.Session) error
	SaveUserData(ctx context.Context, userhash userentity.Hash, userData sessionentity.UserData) error
	Get(ctx context.Context, userhash userentity.Hash, sessionid string) (sessionentity.Session, error)
	Delete(ctx context.Context, userhash userentity.Hash, sessionid string) error
	DeleteAll(ctx context.Context, userhash userentity.Hash) error
}

// New session usecase
func New(repo sessionRepo) *UseCase {
	u := UseCase{
		repo: repo,
	}

	return &u
}

// createID of session
func (u UseCase) createID(userhash userentity.Hash) string {
	keys := []string{time.Now().String(), string(userhash), guuid.New().String()}
	id := strings.Join(keys, ":")
	id = base64.RawStdEncoding.EncodeToString([]byte(id))
	return id
}

// Create function
func (u UseCase) Create(ctx context.Context, userhash userentity.Hash, sessionData sessionentity.Session) (string, error) {
	if err := userhash.Validate(); err != nil {
		return "", err
	}

	id := u.createID(userhash)
	sessionData.ID = id
	err := u.repo.Save(ctx, userhash, id, sessionData)
	if err != nil {
		return "", sessionentity.ErrSessionNotFound
	}

	return id, nil
}

// SetUserInfo to set/change user information in session
func (u UseCase) SetUserInfo(ctx context.Context, userhash userentity.Hash, sessionid string, userData sessionentity.UserData) error {
	sess, err := u.Get(ctx, userhash, sessionid)
	if err != nil {
		return err
	}

	user := userentity.User{}
	bio := userentity.Bio{}

	if userData.User != user {
		sess.UserData.User = userData.User
	}

	if userData.Bio != bio {
		sess.UserData.Bio = userData.Bio
	}

	return nil
}

// Get session
// sessionkey is the user hash-id
// sessionid is the sessionid returned when logged in or creating session
func (u UseCase) Get(ctx context.Context, userhash userentity.Hash, sessionid string) (sessionentity.Session, error) {
	var err error
	s := sessionentity.Session{}

	if userhash == "" {
		return s, sessionentity.Err
	}

	if err := userhash.Validate(); err != nil {
		return s, err
	}

	s, err = u.repo.Get(ctx, userhash, sessionid)
	if err != nil {
		return s, err
	}
	return s, nil
}

// Remove a specific session in a user
func (u UseCase) Remove(ctx context.Context, userhash userentity.Hash, sessionid string) error {
	if err := userhash.Validate(); err != nil {
		return err
	}

	err := u.repo.Delete(ctx, userhash, sessionid)
	if err != nil {
		return err
	}
	return nil
}

// RemoveAll the session in user
func (u UseCase) RemoveAll(ctx context.Context, userhash userentity.Hash) error {
	if err := userhash.Validate(); err != nil {
		return err
	}

	err := u.repo.DeleteAll(ctx, userhash)
	if err != nil {
		return err
	}
	return nil
}
