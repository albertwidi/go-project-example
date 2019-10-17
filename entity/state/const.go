package state

import "time"

// default state expiry time
const (
	DefaultStateExpiryTime = time.Minute * 5
	MinStateExpiryTime     = time.Minute * 1
	MaxStateExpiryTime     = time.Minute * 30
)
