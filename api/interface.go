package api

import (
	"context"

	"github.com/albertwidi/kothak/service/order"
	"github.com/albertwidi/kothak/service/payment"
)

type UserService interface {
	IsUserActive(context.Context, int64) (bool, error)
}

type PaymentService interface {
	CreatePayment(context.Context, payment.Payment) error
	ConfirmPayment(ctx context.Context, paymentid string) error
}

type OrderService interface {
	CreateOrder(context.Context, order.Order) error
	ConfirmOrder(ctx context.Context, orderid string) error
}
