package booking

import (
	"time"
)

// Type of booking
type Type int

// Validate booking type
func (t Type) Validate() error {
	switch t {
	case TypeMonthly:
	case TypeDaily:
	default:
		return ErrInvalidType
	}

	return nil
}

// Status of booking
type Status int

// Booking record attempt of booking
// data in attempt is dirty
// because its mixed from attempt to booking
type Booking struct {
	ID          string
	PropertyID  string
	BookingType Type
	Price       int64
	Deposit     int64
	Status      Status
	// mark that booking really happens
	IsPaid bool
	// a hack, so owner can see booking attempt
	IsCreatedByOwner bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CreatedBy        int64
	IsTest           bool
}

// Detail of booking
type Detail struct {
}

// PaidBooking data
// booking is different attempt
// this is where data after booking confirmation stored
type PaidBooking struct {
	// is an id from attempt
	ID        string
	ItemID    string
	Price     int32
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy int64
	IsTest    bool
}
