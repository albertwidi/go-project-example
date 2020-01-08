package order

// Usecase order
type Usecase struct {
}

// New order usecase
func New() *Usecase {
	u := Usecase{}
	return &u
}
