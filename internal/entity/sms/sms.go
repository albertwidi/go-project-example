package sms

import (
	"errors"

	notifentity "github.com/albertwidi/go-project-example/internal/entity/notification"
)

// Payload struct
type Payload struct {
	From    string
	To      string
	Message string
	Purpose notifentity.Purpose
}

// Validate sms payload
func (p Payload) Validate() error {
	if p.Purpose == notifentity.PurposeEmpty {
		return errors.New("sms: purpose payload cannot be empty")
	}

	return nil
}

// NexmoCallback struct
type NexmoCallback struct {
}
