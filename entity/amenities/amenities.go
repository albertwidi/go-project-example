package amenities

import "time"

// Type of amenities
type Type int

// Amenities struct
type Amenities struct {
	ID        string
	Name      string
	Type      Type
	ImagePath string
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy int64
	IsTest    bool
	IsDeleted bool
}
