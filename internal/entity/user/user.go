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

// // GenderFromString function
// func GenderFromString(gender string) (Gender, error) {
// 	switch gender {
// 	case GenderMaleString:
// 		return GenderMale, nil
// 	case GenderFemaleString:
// 		return GenderFemale, nil
// 	default:
// 		return GenderInvalid, errors.New("gender invalid")
// 	}
// }

// // GenderToString function
// func GenderToString(gender Gender) (string, error) {
// 	switch gender {
// 	case GenderMale:
// 		return GenderMaleString, nil
// 	case GenderFemale:
// 		return GenderFemaleString, nil
// 	default:
// 		return GenderInvalidString, errors.New("gender invalid")
// 	}
// }

// Bio of user
type Bio struct {
	UserID   string
	FullName string
	// Gender     Gender
	Occupation string
	Avatar     string
	Birthday   time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsTest     bool
}
