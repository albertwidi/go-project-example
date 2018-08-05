package payment

type Payment struct {
	ID             string
	IdempotencyKey string
}
