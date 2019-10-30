package session

import (
	"time"

	authentity "github.com/kosanapp/kosan-backend/entity/auth"
	deviceentity "github.com/kosanapp/kosan-backend/entity/device"
	userentity "github.com/kosanapp/kosan-backend/entity/user"
)

// Session struct
type Session struct {
	ID            string              `json:"id"`
	HashID        string              `json:"hash_id"`
	UserData      UserData            `json:"user_data"`
	AuthData      AuthData            `json:"auth_data"`
	KosData       KosData             `json:"kos_data"`
	Device        deviceentity.Device `json:"device_data"`
	Authenticated bool                `json:"authenticated"`
	Tracker       Tracker             `json:"tracker"`
	ExpiryTime    time.Duration       `json:"expiry_time"`
	ExpiredAt     time.Time           `json:"expired_at"`
	CreatedAt     time.Time           `json:"created_at"`
}

// UserData session is a user data cache that assosiated to session of user
type UserData struct {
	User userentity.User `json:"user"`
	Bio  userentity.Bio  `json:"bio"`
}

// AuthData struct
type AuthData struct {
	Provider authentity.Provider `json:"provider"`
}

// KosData struct
type KosData struct {
	// TODO: re-evaluate this value
	ActiveKosID int64   `json:"active_kos_id"`
	KosIDs      []int64 `json:"kos_ids"`
}

// Tracker struct
type Tracker struct {
	LastVisitedPage string `json:"last_visited_page"`
	LastAction      string `json:"last_action"`
}
