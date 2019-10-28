package property

import "time"

// Status of property
type Status int

// Type of property
type Type int

// Segment of the property
type Segment int

// Property data
type Property struct {
	ID      string
	Owner   int64
	Name    string
	Type    Type
	Segment Segment
	// it means we can book the room instead of property
	HasRooms  bool
	Status    Status
	Detail    Detail
	CreatedAt time.Time
	UpdatedAt time.Time
	IsTest    bool
	IsDeleted bool
}

// Detail of property data
type Detail struct {
	PropertyID string
	Address    string
	Amenities  []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsTest     bool
	IsDeleted  bool
}

// Address of property
type Address struct {
	Address   string
	City      string
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsTest    bool
	IsDeleted bool
}

// AddressMap of the property
type AddressMap struct {
	PropertyID string
	Radius     int64
	Lat        float64
	Long       float64
	CreatedAt  time.Time
	IsTest     bool
	IsDeleted  bool
}

// Pricing of property
type Pricing struct {
	PropertyID string
	CreatedAt  time.Time
	IsTest     bool
	IsDeleted  bool
}

// Amenities of property
type Amenities struct {
	PropertyID  string
	AmenitiesID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsTest      bool
	IsDeleted   bool
}
