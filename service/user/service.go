package user

// Resource of user
type Resource interface {
	GetUser(int64) (User, error)
}

// Service of user
type Service struct {
	resource Resource
}

// New user service
func New(userResource Resource) *Service {
	s := Service{
		resource: userResource,
	}
	return &s
}
