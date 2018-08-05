package payment

import (
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
