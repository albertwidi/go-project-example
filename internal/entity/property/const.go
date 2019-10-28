package property

// staus of property
const (
	StatusCreated           = 10
	StatusInactive          = 50
	StatusPermanentlyClosed = 100
	StatusClosed            = 150
	StatusActive            = 200
)

// type list of property
const (
	TypeKos         Type = 1
	TypeHostel      Type = 2
	TypeHotel       Type = 3
	TypeHouse       Type = 4
	TypeApartment   Type = 5
	TypeFlat        Type = 6
	TypePrivateRoom Type = 7
)

// segment list of property
const (
	SegmentStay Segment = 1
)
