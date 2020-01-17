package repository

// Repositories contains list of repository that can be used
type Repositories struct {
}

// New repositories
func New() *Repositories {
	r := Repositories{}
	return &r
}
