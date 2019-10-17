package authentication

import (
	"context"

	authentity "github.com/albertwidi/kothak/entity/authentication"
	stateentity "github.com/albertwidi/kothak/entity/state"
	"github.com/albertwidi/kothak/lib/xerrors"
)

// Usecase of authentication
type Usecase struct {
	stateUsecase stateUsecase
	otpUsecase   otpUsecase
}

type stateUsecase interface {
	Create() error
	Get(ctx context.Context, id string) (stateentity.State, error)
}

type otpUsecase interface {
	Create() error
	Get(ctx context.Context) error
}

// New authentication usecase
func New(stateUsecase stateUsecase, otpUsecase otpUsecase) *Usecase {
	u := Usecase{
		stateUsecase: stateUsecase,
		otpUsecase:   otpUsecase,
	}
	return &u
}

// Try to authenticate
func (u *Usecase) Try(ctx context.Context, username string, action authentity.Action, provider authentity.Provider, metadata map[string]string) error {
	return nil
}

// Authenticate for authenticating user request
func (u *Usecase) Authenticate(ctx context.Context, username, password, stateID string) error {
	op := xerrors.Op("authentication/authenticate")
	_, err := u.stateUsecase.Get(ctx, stateID)
	if err != nil {
		return xerrors.New(op, err)
	}

	return nil
}

// ResendAuthenticationCode for resend the code that used for authentication
// for example the OTP code
func (u *Usecase) ResendAuthenticationCode(ctx context.Context, stateID string) error {
	return nil
}
