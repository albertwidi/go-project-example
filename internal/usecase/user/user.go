package user

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
