package invoice

// list of invoice status
const (
	StatusCreated                       Status = 10
	StatusWaitingForPayment             Status = 100
	StatusWaitingForPaymentConfirmation Status = 105
	StatusPaid                          Status = 200
	StatusCancelled                     Status = 500
	StatusCancelledBySystem             Status = 505
)
