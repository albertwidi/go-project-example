package user

import "errors"

// list of user errors
var (
	ErrUserNotFound = errors.New("user not found")
)
