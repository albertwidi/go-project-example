package user

import (
	"context"

	userentity "github.com/albertwidi/go-project-example/internal/entity/user"
)

// Usecase of user
type Usecase struct {
	authUsecase authUsecase
}

type authUsecase interface {
}

type userRepository interface {
	Create() error
	Update() error
	UpdateStatus() error
	Remove() error
}

// New user usecase
func New(authUsecase authUsecase) *Usecase {
	u := Usecase{
		authUsecase: authUsecase,
	}
	return &u
}

// RegisterData for registration data parameter
type RegisterData struct {
	Country     userentity.Country
	PhoneNumber string
	FullName    string
}

// Validate register data
func (rgd *RegisterData) Validate() error {
	if err := rgd.Country.Validate(); err != nil {
		return err
	}
	return nil
}

// Register usecase
func (u *Usecase) Register(ctx context.Context, data RegisterData) error {
	if err := data.Validate(); err != nil {
		return err
	}
	return nil
}

// RegisterConfirm usecase
func (u *Usecase) RegisterConfirm(ctx context.Context) error {
	return nil
}

// Login usecase
func (u *Usecase) Login(ctx context.Context) error {
	return nil
}
