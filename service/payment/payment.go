package payment

import (
	"context"
	"errors"
)

type Payment struct {
	ID             string
	OrderID        string
	IdempotencyKey string
}

func (s *Service) CreatePayment(ctx context.Context, payment Payment) error {
	if payment.IdempotencyKey == "" {
		return errors.New("payment idempotency key cannot be empty")
	}
	return nil
}

func (s *Service) ConfirmPayment(ctx context.Context, paymentid string) error {
	if paymentid == "" {
		return errors.New("payment id is empty")
	}
	return nil
}
