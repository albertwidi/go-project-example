package redis

import "errors"

// error list
var (
	ErrResponseNotOK = errors.New("redis: response is not ok")
)
