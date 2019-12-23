package notification

import (
	"time"

	"github.com/albertwidi/go-project-example/internal/entity/device"
)

// Client interface
type Client interface {
	Send(templateName string) error
}

// SendOptions struct
type SendOptions struct {
	DryRun bool
}

// UserNotification struct
type UserNotification struct {
	ID             string
	UserID         int64
	ProviderType   int
	ProviderID     int
	ProviderSendID string
	Purpose        int
	IsWebpage      bool
	Status         int
	Title          string
	Message        string
	Show           bool
	HasDetail      bool
	Read           bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsDeleted      bool
	IsTest         bool
	Detail         UserNotificationDetail
}

// UserNotificationDetail struct
type UserNotificationDetail struct {
	NotificationID string
	Body           string
	WebLink        string
	CreatedAt      time.Time
	IsDeleted      bool
	IsTest         bool
}

// Notification struct
type Notification struct {
	UserID     int64
	Title      string
	Message    string
	DetailBody string
	Image      string
	WebLink    string
	Purpose    int
	DeviceID   device.Device
	// either push, sms or email
	NotifData interface{}
}

// Validate notification param
func (n Notification) Validate() error {
	return nil
}

// Options of notification
type Options struct {
	InboxSave bool
	Fake      bool
}

// PushNotification struct
type PushNotification struct {
	// if device token is present, then the device token is used rather than seeking in session
	// this is useful in user register flow
	DeviceToken string
}

// SMSNotification struct
type SMSNotification struct {
	From     string
	ToMSISDN string
}

// EmailNotification struct
type EmailNotification struct {
	To string
	Cc string
}
