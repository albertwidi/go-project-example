package payment

import (
	"context"

	"github.com/albertwidi/kothak/service/payment"
	"github.com/jmoiron/sqlx"
)

type Resource struct {
	masterDB   *sqlx.DB
	followerDB *sqlx.DB
}

func New(masterDB, followerDB *sqlx.DB) *Resource {
	r := Resource{
		masterDB:   masterDB,
		followerDB: followerDB,
	}
	return &r
}

func (r *Resource) CreatePayment(ctx context.Context, pym payment.Payment) error {
	return nil
}
