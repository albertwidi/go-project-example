package order

// status list of order
const (
	// order is updateable
	StatusCreated                       Status = 10
	StatusWaitingForPayment             Status = 100
	StatusWaitingForPaymentConfirmation Status = 105
	StatusPaid                          Status = 150
	StatusWaitingForThirdParty          Status = 180
	StatusComplete                      Status = 200
	StatusRefundInProcess               Status = 400
	// order is not updateable
	StatusOrderRefunded     Status = 450
	StatusCancelled         Status = 500
	StatusCancelledBySystem Status = 505
)

// type list of order
const (
	TypeBooking      Type = 1
	TypeDigitalGoods Type = 2
)

// refund status list
const (
	RefundStatusCreated  RefundStatus = 100
	RefundStatusComplete RefundStatus = 200
)
