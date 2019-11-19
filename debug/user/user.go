package user

import "context"

// DebugUsecase for user
type DebugUsecase struct {
}

// New user debug usecase
func New() *DebugUsecase {
	du := DebugUsecase{}
	return &du
}

// BypassLogin for bypassing log in to the project
func (du *DebugUsecase) BypassLogin(ctx context.Context) error {
	return nil
}
