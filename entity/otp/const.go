package otp

import "time"

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
