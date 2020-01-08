package room

// Usecase room
type Usecase struct {
}

// New room usecase
func New() *Usecase {
	u := Usecase{}
	return &u
}
