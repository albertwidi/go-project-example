package booking

// booking status list
const (
	StatusCreated                       = 10
	StatusWaitingForPayment             = 100
	StatusWaitingForPaymentConfirmation = 105
	StatusWaitingForOwnerConfirmation   = 110
	StatusConfirmed                     = 120
	StatusActive                        = 150
	StatusFinished                      = 200
	StatusCancelled                     = 500
	StatusCancelledBySystem             = 505
	StatusRejected                      = 510
)

// booking type list
const (
	TypeDaily   Type = 1
	TypeMonthly Type = 2
)
