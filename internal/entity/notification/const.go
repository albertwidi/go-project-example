package notification

// type of notification
var (
	TypeSMS         = 1
	TypeEmail       = 2
	TypePushMessage = 3

	TemplateTypePushMessage = "pushmessage"
	TemplateTypeSMS         = "sms"
	TemplateTypeEmail       = "email"
)

// notification providers
var (
	ProviderNexmo = 1
)

// Purpose of notification
type Purpose int

// purpose of notification
var (
	PurposeEmpty                  Purpose = 0
	PurposeSystemUpdate           Purpose = 1
	PurposeAuthenticationOTP      Purpose = 2
	PurposeAuthenticationPayment  Purpose = 3
	PurposeAuthenticationWithdraw Purpose = 4
	PurposePromotion              Purpose = 5
	PurposeReminder               Purpose = 6
)
