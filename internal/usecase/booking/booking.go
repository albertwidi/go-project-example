package booking

// Usecase for booking
type Usecase struct {
}

// New booking usecase
func New() *Usecase {
	u := Usecase{}
	return &u
}

// Create booking
func (u *Usecase) Create() error {
	return nil
}
