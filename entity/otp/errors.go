package otp

import (
	"errors"
)

// list of otp errors
var (
	ErrCodeLengthInvalid          = errors.New("otp: code length is invalid")
	ErrOTPInvalid                 = errors.New("otp: password invalid")
	ErrOTPExpired                 = errors.New("otp: otp already expired")
	ErrOTPReachResendMaxAttempt   = errors.New("otp: reach maximum resend attempt")
	ErrOTPReachValidateMaxAttempt = errors.New("otp: reach maximum validate attempt, too many wrong password")
	ErrOTPNotResendable           = errors.New("otp: not resendable")
)
