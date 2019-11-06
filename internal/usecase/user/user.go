package user

import (
	"context"

	userentity "github.com/albertwidi/go_project_example/internal/entity/user"
)

// Usecase of user
type Usecase struct {
	authUsecase authUsecase
}

type authUsecase interface {
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

// Register usecase
func (u *Usecase) Register(ctx context.Context, data RegisterData) error {
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
