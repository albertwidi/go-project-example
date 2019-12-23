package session

import (
	"time"

	authentity "github.com/albertwidi/go-project-example/internal/entity/authentication"
	userentity "github.com/albertwidi/go-project-example/internal/entity/user"
)

// Session struct
type Session struct {
	ID            string        `json:"id"`
	HashID        string        `json:"hash_id"`
	AuthData      AuthData      `json:"auth_data"`
	Authenticated bool          `json:"authenticated"`
	ExpiryTime    time.Duration `json:"expiry_time"`
	ExpiredAt     time.Time     `json:"expired_at"`
	CreatedAt     time.Time     `json:"created_at"`
}

// UserData session is a user data cache that assosiated to session of user
type UserData struct {
	User userentity.User `json:"user"`
	Bio  userentity.Bio  `json:"bio"`
}

// AuthData struct
type AuthData struct {
	Provider authentity.Provider `json:"provider"`
	Action   authentity.Action   `json:"action"`
}
