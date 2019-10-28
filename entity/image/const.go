package image

// list of image mode
const (
	ModePublic  Mode = "public"
	ModePrivate Mode = "private"
	ModeSigned  Mode = "signed"
)

// list of group of image
const (
	GroupEmpty          = ""
	GroupMixed          = "mixed"
	GroupPropertyKos    = "property/kos"
	GroupPropertyHotel  = "property/hotel"
	GroupPropertyHostel = "property/hostel"
	GroupPropertyHouse  = "property/house"
	GroupPropertyRoom   = "property/room"
	GroupAmenities      = "amenities"
	GroupUserKTP        = "user/ktp"
	GroupUserAvatar     = "user/avatar"
	GroupPaymentProof   = "payment/proof"
)
