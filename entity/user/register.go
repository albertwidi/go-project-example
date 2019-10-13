package user

import "time"

// Registrations struct
type Registrations struct {
	ID          int64
	UserID      int64
	UserType    int
	UserStatus  int
	HashID      string
	KTPID       int64
	FullName    string
	Email       string
	PhoneNumber string
	Gender      int
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
