package payment

import "context"

type Resource interface {
	CreatePayment(context.Context, Payment) error
}

type UserService interface {
	IsUserActive(context.Context, int64) (bool, error)
}

type Service struct {
	resource Resource
	userSvc  UserService
}

// New payment service
func New(paymentResource Resource, userService UserService) *Service {
	s := Service{
		resource: paymentResource,
		userSvc:  userService,
	}
	return &s
}
