package authentication

import (
	"context"
	"strings"
	"time"

	authentity "github.com/albertwidi/go-project-example/internal/entity/authentication"
	otpentity "github.com/albertwidi/go-project-example/internal/entity/otp"
	stateentity "github.com/albertwidi/go-project-example/internal/entity/state"
	"github.com/albertwidi/go-project-example/internal/xerrors"
)

// Usecase of authentication
type Usecase struct {
	stateUsecase stateUsecase
	otpUsecase   otpUsecase
}

type stateUsecase interface {
	Create(ctx context.Context, state stateentity.State) (string, error)
	Get(ctx context.Context, id string) (stateentity.State, error)
}

type secretUsecase interface {
	Get(ctx context.Context, username string) error
}

type otpUsecase interface {
	Create(ctx context.Context, uniqueID string, codeLength otpentity.CodeLength, expire time.Duration) (*otpentity.OTP, error)
	Get(ctx context.Context, uniqueID string) error
}

// New authentication usecase
func New(stateUsecase stateUsecase, otpUsecase otpUsecase) *Usecase {
	u := Usecase{
		stateUsecase: stateUsecase,
		otpUsecase:   otpUsecase,
	}
	return &u
}

// Authenticate for trying to authenticate
func (u *Usecase) Authenticate(ctx context.Context, username string, action authentity.Action, provider authentity.Provider, metadata map[string]string) (string, error) {
	state := stateentity.State{
		CreatedBy: username,
		Authentication: authentity.Authentication{
			Action:   action,
			Provider: provider,
			Username: username,
		},
		MetaData:  metadata,
		CreatedAt: time.Now(),
	}

	id, err := u.stateUsecase.Create(ctx, state)
	if err != nil {
		return "", err
	}
	switch provider {
	case authentity.ProviderOTP:
		otpUniqueID := strings.Join([]string{username, string(action)}, ",")
		_, err := u.otpUsecase.Create(ctx, otpUniqueID, otpentity.CodeLength6, time.Minute*5)
		if err != nil {
			return "", err
		}

	}
	return id, nil
}

// Confirm for confirming the authenticating user request
func (u *Usecase) Confirm(ctx context.Context, username, password, stateID string) error {
	op := xerrors.Op("authentication/authenticate")
	_, err := u.stateUsecase.Get(ctx, stateID)
	if err != nil {
		return xerrors.New(op, err)
	}
	return nil
}

// ResendCode for resend the code that used for authentication
// for example the OTP code
func (u *Usecase) ResendCode(ctx context.Context, stateID string) error {
	return nil
}

// IsAuthenticated to check whether state is used and authenticated for particular username
func (u *Usecase) IsAuthenticated(ctx context.Context, username, stateID string, action authentity.Action) (bool, error) {
	return false, nil
}
