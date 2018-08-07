package order

import (
	"context"
	"errors"
)

type Order struct {
	ID             string
	UserID         int64
	IdempotencyKey string
	Metadata       string
}

func (s *Service) CreateOrder(ctx context.Context, order Order) error {
	if order.IdempotencyKey == "" {
		return errors.New("order idempotency key cannot be empty")
	}
	// check wether user is active or not using user service
	if userActive, err := s.userSvc.IsUserActive(ctx, order.UserID); err != nil {
		return err
	} else if !userActive {
		return errors.New("cannot create order, user is not active")
	}

	return s.resource.CreateOrder(ctx, order)
}

func (s *Service) ConfirmOrder(ctx context.Context, orderid string) error {
	if orderid == "" {
		return errors.New("orderid is empty")
	}
	return nil
}
