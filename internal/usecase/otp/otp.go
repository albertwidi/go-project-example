package otp

import (
	"context"
	"strconv"
	"time"

	otpentity "github.com/albertwidi/go-project-example/internal/entity/otp"
	"github.com/albertwidi/go-project-example/internal/pkg/randgen"
	"github.com/albertwidi/go-project-example/internal/xerrors"
)

// Usecase of otp
type Usecase struct {
	repo otpRepo
	// code generator
	rgen map[otpentity.CodeLength]*randgen.Generator
}

type otpRepo interface {
	Save(ctx context.Context, otp otpentity.OTP) error
	SetLast(ctx context.Context, otp otpentity.OTP) error
	GetLast(ctx context.Context, uniqueID string) (otpentity.OTP, error)
	Len(ctx context.Context, uniqueID string) (int, error)
	IncreaseValidateAttempt(ctx context.Context, uniqueID string, attempt int) (int, error)
	DeleteValidateAttempt(ctx context.Context, uniqueID string) error
	DeleteAll(ctx context.Context, uniqueID string) error
}

// New otp usecase
func New(otpRepo otpRepo) (*Usecase, error) {
	workerNumber := 3
	rgen := make(map[otpentity.CodeLength]*randgen.Generator)
	rgen[otpentity.CodeLength4] = randgen.New(workerNumber, 1000, 9999, time.Now().UnixNano())
	rgen[otpentity.CodeLength6] = randgen.New(workerNumber, 100000, 999999, time.Now().UnixNano())

	u := Usecase{
		repo: otpRepo,
		rgen: rgen,
	}
	return &u, nil
}

// Create OTP
func (u Usecase) Create(ctx context.Context, uniqueID string, codeLength otpentity.CodeLength, expire time.Duration) (*otpentity.OTP, error) {
	if err := codeLength.Validate(); err != nil {
		return nil, err
	}

	var (
		otpcode        = strconv.Itoa(u.rgen[codeLength].Generate())
		resendDuration = otpentity.ResendTimeDefault
	)

	lastOTP, err := u.repo.GetLast(ctx, uniqueID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if err != nil {
		return nil, err
	}

	// means the last otp is not exists
	if lastOTP.Code != "" {
		// make sure that otp is not resendable before the duration
		if now.Before(lastOTP.ResendableAt) {
			return nil, xerrors.New(otpentity.ErrOTPNotResendable)
		}

		// get the length of current otp
		otpLen, err := u.repo.Len(ctx, uniqueID)
		if err != nil {
			return nil, err
		}

		// means otp is on its threshold
		if otpLen%otpentity.ThresholdOTPResend == 0 {
			resendDuration = otpentity.ResendTimeAfterThreshold
		}
	}

	if expire == 0 {
		expire = otpentity.ExpiryTimeDefault
	}

	otp := otpentity.OTP{
		UniqueID:     uniqueID,
		Code:         otpcode,
		CreatedAt:    now,
		ExpiryTime:   expire,
		ExpiredAt:    now.Add(expire),
		ResendTime:   otpentity.ResendTimeDefault,
		ResendableAt: now.Add(resendDuration),
	}

	err = u.repo.Save(ctx, otp)
	if err != nil {
		return nil, err
	}

	return &otp, nil
}

// Validate the otpcode
func (u Usecase) Validate(ctx context.Context, uniqueID string, otpCode string) error {
	// always get the last otp from certain uniqueid and action
	// this is because otp is grouped per uniqueid and action
	otp, err := u.repo.GetLast(ctx, uniqueID)
	if err != nil {
		return err
	}

	now := time.Now()
	if err != nil {
		return err
	}

	// check whether the otp already expired
	if now.After(otp.ExpiredAt) {
		return xerrors.New(otpentity.ErrOTPExpired)
	}

	// check whether the otp is valid
	if otp.Code != otpCode {
		validateAttempt, err := u.repo.IncreaseValidateAttempt(ctx, uniqueID, 1)
		if err != nil {
			return err
		}

		// if otp validation reach its threshold, invalidate and expire otp immediately
		if validateAttempt%otpentity.ThresholdOTPValidate == 0 {
			otp.ExpiryTime = 0
			otp.ExpiredAt = now
			otp.ResendTime = otpentity.ResendTimeMax
			otp.ResendableAt = now.Add(otpentity.ResendTimeMax)
		}

		if err := u.repo.SetLast(ctx, otp); err != nil {
			return err
		}

		return xerrors.New(otpentity.ErrOTPInvalid)
	}

	err = u.repo.DeleteValidateAttempt(ctx, uniqueID)
	if err != nil {
		return err
	}

	return nil
}

// Delete OTP
func (u Usecase) Delete(ctx context.Context, uniqueID string) error {
	err := u.repo.DeleteAll(ctx, uniqueID)
	if err != nil {
		return err
	}

	return err
}
