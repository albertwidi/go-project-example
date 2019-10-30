package image

import (
	"fmt"
	"net/textproto"

	userentity "github.com/albertwidi/go_project_example/internal/entity/user"
)

// Mode of image
type Mode string

// Validate mode
func (m Mode) Validate() error {
	switch m {
	case ModePrivate:
		break
	case ModePublic:
		break
	case ModeSigned:
		break
	default:
		return fmt.Errorf("image: invalid mode, got %s", m)
	}

	return nil
}

// Group of image
type Group string

// Validate group
func (g Group) Validate() error {
	switch g {
	case GroupAmenities,
		GroupPropertyKos,
		GroupPropertyRoom,
		GroupPropertyHotel,
		GroupPropertyHostel,
		GroupPropertyHouse,
		GroupPaymentProof,
		GroupUserKTP,
		GroupUserAvatar:
	default:
		return fmt.Errorf("image: invalid group, got %s", g)
	}

	return nil
}

// FileInfo struct
type FileInfo struct {
	FileName string
	Size     int64
	Header   textproto.MIMEHeader
	UserHash userentity.Hash
	Mode     Mode
	Group    Group
	Tags     string
}

// Options struct
type Options struct {
	Manipulation *Manipulation
}

// Manipulation struct
type Manipulation struct {
	Resize Resize
}

// Resize struct
type Resize struct {
	Height int
	Width  int
}
