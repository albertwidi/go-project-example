package user

import "time"

// Registrations struct
type Registrations struct {
	ID          string
	UserID      string
	UserType    int
	UserStatus  int
	HashID      string
	FullName    string
	Email       string
	PhoneNumber string
	BirthDate   time.Time
	Channel     int
	Device      int
	Lat         string
	Long        string
	DeviceToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsTest      bool
}
