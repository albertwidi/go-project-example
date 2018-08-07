package order

import "context"

type Resource interface {
	CreateOrder(context.Context, Order) error
}

type UserService interface {
	IsUserActive(context.Context, int64) (bool, error)
}

// Service order
type Service struct {
	resource Resource
	userSvc  UserService
}

// New order service
func New(orderResource Resource, userService UserService) *Service {
	s := Service{
		resource: orderResource,
		userSvc:  userService,
	}
	return &s
}
