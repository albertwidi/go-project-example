package user

import (
	"errors"
	"fmt"
	"time"
)

// list of user types
type (
	// Type of user
	Type int
	// Hash of user
	Hash string
	// Country of user
	Country string
	// CountryCode of user
	CountryCode string
	// PhoneNumber of user
	PhoneNumber string
)

// Validate user hash
func (h Hash) Validate() error {
	if h == "" {
		return errors.New("user hash cannot be empty")
	}

	return nil
}

// Validate user country
func (c Country) Validate() error {
	switch c {
	case CountryID:
		return nil
	}
	return fmt.Errorf("country: country with name %s is not valid", c)
}

// User struct
type User struct {
	ID          string
	HashID      Hash
	UserStatus  int
	UserType    Type
	PhoneNumber string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsTest      bool
}

// Bio of user
type Bio struct {
	UserID   string
	FullName string
	// Gender     Gender
	Avatar    string
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	IsTest    bool
}
