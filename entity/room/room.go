package room

// Room struct
type Room struct {
	ID   string
	Tags []string
	Name string
}

// Detail of the room
type Detail struct {
	RoomID    string
	Amenities []string
}

// PricingType of room
type PricingType int

// Pricing of room
type Pricing struct {
	Name   string
	Amount string
}

// Amenities of room
type Amenities struct {
	RoomID      string
	AmenitiesID string
}
