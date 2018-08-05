package order

import (
	"context"

	"github.com/jmoiron/sqlx"
	orderservice "gitlab.com/kosanapp/kothak/services/order"
)

type Resource struct {
	masterDB   *sqlx.DB
	followerDB *sqlx.DB
}

// New order resource
func New(masterDB, followerDB *sqlx.DB) *Resource {
	r := Resource{
		masterDB:   masterDB,
		followerDB: followerDB,
	}
	return &r
}

func (r *Resource) CreateOrder(ctx context.Context, order orderservice.Order) error {
	return nil
}
