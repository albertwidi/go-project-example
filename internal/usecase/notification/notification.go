package notification

import (
	"context"
	"fmt"

	smsentity "github.com/albertwidi/go-project-example/internal/entity/sms"
	pushentity "github.com/albertwidi/go-project-example/internal/entity/pushmessage"
	entity "github.com/albertwidi/go-project-example/internal/entity/notification"
)

// UseCase struct
type UseCase struct {
	repo     notificationRepo
	pushrepo pushMessageRepo
	smsrepo  smsRepo
}

type notificationRepo interface {
	Save(ctx context.Context, data entity.UserNotification) error
	Get(ctx context.Context, userID int64) ([]entity.UserNotification, error)
	GetDetail(ctx context.Context, userID int64, notifID string) (entity.UserNotificationDetail, error)
}

type pushMessageRepo interface {
	Send(ctx context.Context, message pushentity.Message) error
}

type smsRepo interface {
	Send(ctx context.Context, payload smsentity.Payload) error
}

// New usecase of notification
func New(repo notificationRepo, pushrepo pushMessageRepo, smsrepo smsRepo) *UseCase {
	u := UseCase{
		repo:     repo,
		pushrepo: pushrepo,
		smsrepo:  smsrepo,
	}

	return &u
}

// Send notification to user
func (u *UseCase) Send(ctx context.Context, notification entity.Notification, options *entity.Options) error {
	// validate the request first
	if err := notification.Validate(); err != nil {
		return err
	}

	var (
		notificationType int
		opts             entity.Options
		err              error
	)

	if options != nil {
		opts = *options
	}

	userNotifDetail := entity.UserNotificationDetail{
		Body:    notification.DetailBody,
		WebLink: notification.WebLink,
	}

	userNotif := entity.UserNotification{
		UserID:  notification.UserID,
		Title:   notification.Title,
		Message: notification.Message,
		Detail:  userNotifDetail,
	}

	// send the actual notification
	switch t := notification.NotifData.(type) {
	// send push notification
	case entity.PushNotification:
		notificationType = entity.TypePushMessage
		data := notification.NotifData.(entity.PushNotification)

		message := pushentity.Message{
			Token: data.DeviceToken,
			Title: notification.Title,
			Body:  notification.Message,
		}

		err = u.pushrepo.Send(ctx, message)

	case entity.SMSNotification:
		notificationType = entity.TypeSMS
		data := notification.NotifData.(entity.SMSNotification)

		payload := smsentity.Payload{
			From:    data.From,
			To:      data.ToMSISDN,
			Message: notification.Message,
		}

		err = u.smsrepo.Send(ctx, payload)

	case entity.EmailNotification:
		notificationType = entity.TypeEmail

	default:
		return fmt.Errorf("notification: notification data type is not valid, got %v", t)
	}

	if err != nil {
		return err
	}

	userNotif.ProviderID = notificationType
	// save notification
	if opts.InboxSave {
		if err := u.repo.Save(ctx, userNotif); err != nil {
			return err
		}
	}

	return nil
}

// Get user notification
func (u *UseCase) Get(ctx context.Context, userID int64) ([]entity.UserNotification, error) {
	userNotif, err := u.repo.Get(ctx, userID)
	if err != nil {
		return userNotif, err
	}

	return userNotif, nil
}

// GetDetail of notification
func (u *UseCase) GetDetail(ctx context.Context, notifID string, userID int64) (entity.UserNotificationDetail, error) {
	notifDetail, err := u.repo.GetDetail(ctx, userID, notifID)
	if err != nil {
		return notifDetail, err
	}

	return notifDetail, nil
}
