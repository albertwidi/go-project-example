package otp

import "time"

// CodeLength type
type CodeLength int

// Validate code length
func (clength CodeLength) Validate() error {
	if clength < CodeLength4 || clength > CodeLength6 {
		return ErrCodeLengthInvalid
	}
	return nil
}

// legnth of code for otp
const (
	CodeLength4 = 4
	CodeLength6 = 6
)

// list of otp threshold
const (
	ThresholdOTPResend   int = 3
	ThresholdOTPValidate int = 10
)

// list of expiry time
const (
	OTPKeyExpiry = time.Hour * 24

	// expire time for otp
	ExpiryTimeDefault = time.Minute * 4

	// the default time for expiring otp is 90 seconds
	ResendTimeDefault        = time.Second * 90
	ResendTimeAfterThreshold = time.Minute * 30
	// the maximum time for otp exiry is 1 day
	ResendTimeMax = time.Hour * 24
)
