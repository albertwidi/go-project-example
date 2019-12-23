package otp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	authentity "github.com/albertwidi/go-project-example/internal/entity/authentication"
	otpentity "github.com/albertwidi/go-project-example/internal/entity/otp"
	"github.com/albertwidi/go-project-example/internal/pkg/conv"
	"github.com/albertwidi/go-project-example/internal/pkg/redis"
)

// Repository for otp
type Repository struct {
	redis redis.Redis
}

// New otp repo
func New(redis redis.Redis) *Repository {
	r := Repository{
		redis: redis,
	}
	return &r
}

// Save the otp
func (r Repository) Save(ctx context.Context, otp otpentity.OTP) error {
	otpKey := generateOtpKey(otp.UniqueID, otp.Action)

	flattenedotp := r.flattenOTP(otp)
	_, err := r.redis.LPush(ctx, otpKey, flattenedotp)
	if err != nil {
		return err
	}

	// expire the key
	// TODO: re-think about the expiration, as of now one expire for all otp
	_, err = r.redis.Expire(ctx, otpKey, int(otpentity.OTPKeyExpiry))
	return err
}

// SetLast value to wanted otp value
func (r Repository) SetLast(ctx context.Context, otp otpentity.OTP) error {
	otpKey := generateOtpKey(otp.UniqueID, otp.Action)
	index := 0

	flattenedotp := r.flattenOTP(otp)
	_, err := r.redis.LSet(ctx, otpKey, flattenedotp, index)
	return err
}

// GetLast OTP
func (r Repository) GetLast(ctx context.Context, uniqueID string, action authentity.Action) (otpentity.OTP, error) {
	index := 0

	otpKey := generateOtpKey(uniqueID, action)
	result, err := r.redis.LIndex(ctx, otpKey, index)
	if err != nil {
		return otpentity.OTP{}, err
	}

	otp := otpentity.OTP{}
	if otp, err = r.deflateOTP(result); err != nil {
		return otp, err
	}

	return otp, nil
}

// Len of the otp
func (r Repository) Len(ctx context.Context, uniqueID string, action authentity.Action) (int, error) {
	otpKey := generateOtpKey(uniqueID, action)
	result, err := r.redis.LLen(ctx, otpKey)
	if err != nil {
		return 0, err
	}

	return result, err
}

// IncreaseValidateAttempt to increase number when people trying to validate their otp
func (r Repository) IncreaseValidateAttempt(ctx context.Context, uniqueID string, amount int) (int, error) {
	valKey := generateOtpValidateAttemptKey(uniqueID)

	result, err := r.redis.IncrementBy(ctx, valKey, amount)
	if err != nil {
		return 0, err
	}

	// set the expiry of validate attempt
	if _, err := r.redis.Expire(ctx, valKey, int(durationExpireVAlidateAttempt)); err != nil {
		return 0, err
	}

	return result, err
}

// DeleteValidateAttempt to delete the attempt of validation
// this usually happened when people finally able to validate their OTP
func (r Repository) DeleteValidateAttempt(ctx context.Context, uniqueID string) error {
	valKey := generateOtpValidateAttemptKey(uniqueID)

	_, err := r.redis.Delete(ctx, valKey)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAll otp from a given user with uniqueID and authentication action
func (r Repository) DeleteAll(ctx context.Context, uniqueID string, action authentity.Action) error {
	otpKey := generateOtpKey(uniqueID, action)
	_, err := r.redis.Delete(ctx, otpKey)
	if err != nil {
		return err
	}

	return nil
}

// flatten the OTP
func (r Repository) flattenOTP(otp otpentity.OTP) string {
	return fmt.Sprintf("%s:%s:%d:%d:%d:%d", otp.UniqueID, otp.Code, otp.ExpiryTime, otp.ExpiredAt.Unix(), otp.ResendTime, otp.ResendableAt.Unix())
}

func (r Repository) deflateOTP(flattenedOTP string) (otpentity.OTP, error) {
	var (
		err error
		otp otpentity.OTP
	)

	val := strings.Split(flattenedOTP, ":")
	if len(val) != 6 {
		return otp, errors.New("not a valid otp format")
	}

	otp.UniqueID = val[0]
	if err != nil {
		return otp, err
	}
	otp.Code = val[1]

	expTime, err := conv.StringToInt64(val[2])
	if err != nil {
		return otp, err
	}
	otp.ExpiryTime = time.Duration(expTime)

	timeExpire, err := conv.StringToInt64(val[3])
	if err != nil {
		return otp, err
	}
	otp.ExpiredAt = time.Unix(timeExpire, 0)

	resendTime, err := conv.StringToInt64(val[4])
	if err != nil {
		return otp, err
	}
	otp.ResendTime = time.Duration(resendTime)

	timeResend, err := conv.StringToInt64(val[5])
	if err != nil {
		return otp, err
	}
	otp.ResendableAt = time.Unix(timeResend, 0)

	return otp, nil
}

func generateOtpKey(uniqueID string, action authentity.Action) string {
	return fmt.Sprintf(keyOTP, uniqueID, action)
}

func generateOtpValidateAttemptKey(uniqueID string) string {
	return fmt.Sprintf(keyValidateAttempt, uniqueID)
}
