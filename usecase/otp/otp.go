package otp

import (
	"context"
	"strconv"
	"time"

	authentity "github.com/kosanapp/kosan-backend/entity/auth"
	otpentity "github.com/kosanapp/kosan-backend/entity/otp"
	kosanerr "github.com/kosanapp/kosan-backend/kosan/errors"
	"github.com/kosanapp/kosan-backend/pkg/errors"
	"github.com/kosanapp/kosan-backend/pkg/randgen"
	"github.com/kosanapp/kosan-backend/pkg/redis"
	"github.com/kosanapp/kosan-backend/pkg/timeutil"
)

// Usecase of otp
type Usecase struct {
	repo otpRepo
	rgen *randgen.Generator
}

type otpRepo interface {
	Save(ctx context.Context, otp otpentity.OTP) error
	SetLast(ctx context.Context, otp otpentity.OTP) error
	GetLast(ctx context.Context, uniqueID string, action authentity.Action) (otpentity.OTP, error)
	Len(ctx context.Context, uniqueID string, action authentity.Action) (int, error)
	IncreaseValidateAttempt(ctx context.Context, uniqueID string, attempt int) (int, error)
	DeleteValidateAttempt(ctx context.Context, uniqueID string) error
	DeleteAll(ctx context.Context, uniqueID string, action authentity.Action) error
}

// New otp usecase
func New(otpRepo otpRepo) (*Usecase, error) {
	workerNumber := 3
	minOTP := 100000
	maxOTP := 999999

	u := Usecase{
		repo: otpRepo,
		rgen: randgen.New(workerNumber, minOTP, maxOTP, time.Now().UnixNano()),
	}

	return &u, nil
}

// Create OTP
func (u Usecase) Create(ctx context.Context, uniqueID string, action authentity.Action) (otpentity.OTP, error) {
	var (
		otpcode        = strconv.Itoa(u.rgen.Generate())
		resendDuration = otpentity.ResendTimeDefault
	)

	lastOTP, err := u.repo.GetLast(ctx, uniqueID, action)
	if err != nil && !redis.IsErrNil(err) {
		return otpentity.OTP{}, err
	}

	now, err := timeutil.TimeNowWIB()
	if err != nil {
		return otpentity.OTP{}, err
	}

	// means the last otp is not exists
	if lastOTP.Code != "" {
		// make sure that otp is not resendable before the duration
		if now.Before(lastOTP.ResendableAt) {
			return lastOTP, errors.E(kosanerr.ErrOTPSendNeedToWait)
		}

		// get the length of current otp
		otpLen, err := u.repo.Len(ctx, uniqueID, action)
		if err != nil {
			return otpentity.OTP{}, err
		}

		// means otp is on its threshold
		if otpLen%otpentity.ThresholdOTPResend == 0 {
			resendDuration = otpentity.ResendTimeAfterThreshold
		}
	}

	otp := otpentity.OTP{
		UniqueID:     uniqueID,
		Action:       action,
		Code:         otpcode,
		CreatedAt:    now,
		ExpiryTime:   otpentity.ExpiryTimeDefault,
		ExpiredAt:    now.Add(otpentity.ExpiryTimeDefault),
		ResendTime:   otpentity.ResendTimeDefault,
		ResendableAt: now.Add(resendDuration),
	}

	err = u.repo.Save(ctx, otp)
	if err != nil {
		return otp, err
	}

	return otp, nil
}

// Recreate otp using last otp
// TODO: reuse create function
func (u Usecase) Recreate(ctx context.Context, uniqueID string, action authentity.Action) (otpentity.OTP, error) {
	otp, err := u.repo.GetLast(ctx, uniqueID, action)
	if err != nil {
		return otp, err
	}

	ok, err := otp.IsResendable()
	if err != nil {
		return otp, err
	}

	if !ok {
		return otp, err
	}

	otp, err = u.Create(ctx, otp.UniqueID, otp.Action)
	if err != nil {
		return otp, err
	}

	return otp, nil
}

// Validate the otpcode
func (u Usecase) Validate(ctx context.Context, uniqueID string, action authentity.Action, otpCode string) error {
	// always get the last otp from certain uniqueid and action
	// this is because otp is grouped per uniqueid and action
	otp, err := u.repo.GetLast(ctx, uniqueID, action)
	if err != nil {
		return err
	}

	now, err := timeutil.TimeNowWIB()
	if err != nil {
		return err
	}

	// check whether the otp already expired
	if now.After(otp.ExpiredAt) {
		err = errors.E(kosanerr.ErrOTPAlreadyExpired)
		return err
	}

	// check whether the otp is valid
	if otp.Code != otpCode {
		validateAttempt, err := u.repo.IncreaseValidateAttempt(ctx, uniqueID, 1)
		if err != nil {
			return err
		}

		if validateAttempt%otpentity.ThresholdOTPValidate == 0 {
			otp.ExpiryTime = 0
			otp.ExpiredAt = now
			otp.ResendTime = otpentity.ResendTimeMax
			otp.ResendableAt = now.Add(otpentity.ResendTimeMax)
		}

		if err := u.repo.SetLast(ctx, otp); err != nil {
			return err
		}

		err = kosanerr.ErrOTPInvalid
		return err
	}

	err = u.repo.DeleteValidateAttempt(ctx, uniqueID)
	if err != nil {
		return err
	}

	return nil
}

// Delete OTP
func (u Usecase) Delete(ctx context.Context, uniqueID string, action authentity.Action) error {
	err := u.repo.DeleteAll(ctx, uniqueID, action)
	if err != nil {
		return err
	}

	return err
}
