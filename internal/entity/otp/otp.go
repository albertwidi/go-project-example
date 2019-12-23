package otp

import (
	"time"

	authentity "github.com/albertwidi/go-project-example/internal/entity/authentication"
)

// OTP struct
type OTP struct {
	UniqueID     string            `json:"unique_id"`
	Action       authentity.Action `json:"action"`
	Code         string            `json:"code"`
	CreatedAt    time.Time         `json:"created_at"`
	ExpiryTime   time.Duration     `json:"expiry_time"`
	ExpiredAt    time.Time         `json:"expired_at"`
	ResendTime   time.Duration     `json:"resend_time"`
	ResendableAt time.Time         `json:"resendable_at"`
}

// IsResendable return whether otp is resendable or not
func (otp OTP) IsResendable() (bool, error) {
	if otp.Code == "" {
		return true, nil
	}

	now := time.Now()
	if now.Before(otp.ResendableAt) {
		return false, nil
	}

	return true, nil
}
