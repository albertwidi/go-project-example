package state

import (
	"errors"
	"time"

	authentity "github.com/albertwidi/kothak/entity/authentication"
)

// New to create a new state
func New() State {
	s := State{
		Data:     make(map[string]interface{}),
		MetaData: make(map[string]string),
	}

	return s
}

// State data
type State struct {
	// CreatedBy is a validator that state is created by some user
	CreatedBy int64
	// CreatedByHashID is a hash of user id used externally
	CreatedByHashID string
	// identifier can be anything, from phone_number, user_id or unique_id, depends on needs
	// this is useful for authenticating request
	Identifier string
	// state might be used for authentication
	// so we need to record the authentication data
	Authentication authentity.Authentication
	// navigation when using state
	// we might want to redirect or pointing to some pages after state validation
	Navigation Navigation
	// data that might want to be stored in the creation of state
	Data map[string]interface{}
	// metadata that want to be stored in the creation of state
	MetaData   map[string]string
	ExpiryTime time.Duration
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

// Validate state
func (s State) Validate() error {
	if s.CreatedBy == 0 {
		return errors.New("state created by cannot be empty")
	}

	if s.Identifier == "" {
		return errors.New("state identifier cannot be empty")
	}

	if s.ExpiryTime > MaxStateExpiryTime {
		return errors.New("expiry time is more than max expiry time")
	} else if s.ExpiryTime < MinStateExpiryTime {
		return errors.New("expirty time is less than min expiry time")
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

// Navigation of state
type Navigation struct {
	AppRoute string
}
