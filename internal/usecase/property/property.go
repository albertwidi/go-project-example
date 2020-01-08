package property

// Usecase of property
type Usecase struct {
}

// New property usecase
func New() *Usecase {
	u := Usecase{}
	return &u
}
