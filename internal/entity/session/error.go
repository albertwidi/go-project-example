package session

import "errors"

// list of error
var (
	ErrSessionNotFound error = errors.New("session: not found")
)
