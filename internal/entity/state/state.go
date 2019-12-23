package state

import (
	"time"

	authentity "github.com/albertwidi/go-project-example/internal/entity/authentication"
)

// New to create a new state
func New() State {
	s := State{
		MetaData: make(map[string]string),
	}

	return s
}

// State data
type State struct {
	// CreatedBy define by whom the state is created
	CreatedBy string
	// CreatedByHashID is a hash of user id used externally
	CreatedByHashID string
	// Authentication data of state
	Authentication authentity.Authentication
	// metadata that want to be stored in the creation of state
	MetaData   map[string]string
	ExpiryTime time.Duration
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

// Validate state
func (s State) Validate() error {
	if s.CreatedBy == "" {
		return ErrCreatedByEmpty
	}

	if s.ExpiryTime > MaxStateExpiryTime {
		return ErrExpiryTimeMoreThanMax
	} else if s.ExpiryTime < MinStateExpiryTime {
		return ErrExpiryTimeLessThanMin
	}

	return nil
}

// IsExpired to check whether the state is expired or not
func (s State) IsExpired() (bool, error) {
	now := time.Now()

	if now.After(s.ExpiredAt) {
		return true, nil
	}

	return false, nil
}
