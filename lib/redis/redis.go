package redis

import (
	"context"
	"errors"
)

// error list
var (
	ErrResponseNotOK = errors.New("redis: response is not ok")
)

// Redis interface
type Redis interface {
	SetEx(ctx context.Context, key string, value interface{}, expire int) (string, error)
}
