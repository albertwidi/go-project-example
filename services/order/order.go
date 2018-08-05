package order

type Order struct {
	ID             string
	IdempotencyKey string
	Metadata       string
}
