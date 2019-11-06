package secret

import (
	"time"
)

// Key of secret
type Key string

// Secret of user
type Secret struct {
	ID          string
	UserID      string
	SecretKey   Key
	SecretValue string
	CreatedAt   time.Time
	CreatedBy   int64
	UpdatedAt   time.Time
	UpdatedBy   int64
	IsTest      bool
}
