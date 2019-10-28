package state

import "errors"

// error list of state
var (
	ErrStateNotFound         = errors.New("state: not found")
	ErrCreatedByEmpty        = errors.New("state: created by is empty")
	ErrExpiryTimeMoreThanMax = errors.New("state: expiry time is more than max expiry time allowed")
	ErrExpiryTimeLessThanMin = errors.New("state: expiry time is  less than min expiry time allowed")
)
