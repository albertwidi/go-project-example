package invoice

import "time"

// Status of invoice
type Status int

// Invoice struct
type Invoice struct {
	ID            string
	Number        string
	OrderID       string
	InvoiceFrom   string
	InvoiceTo     string
	Type          int
	Total         int64
	DiscountTotal int64
	GrandTotal    int64
	Details       []Detail
	Status        Status
	Description   string
	CreatedBy     int64
	DueDate       time.Time
	PaidAt        time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsTest        bool
	IsDeleted     bool
}

// Detail for invoice detail
type Detail struct {
	ID           string
	InvoiceID    string
	Amount       int64
	Discount     int64
	ItemName     string
	ItemQuantity int64
	Description  string
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	IsTest       bool
	IsDeleted    bool
}

// PaidInvoice struct
// store all data of paid invoice
type PaidInvoice struct {
	InvoiceID     string
	InvoiceNumber string
	PaymentID     string
	OrderID       string
	InvoiceFrom   int64
	InvoiceTo     int64
	CreatedAt     time.Time
	CreatedBy     int64
	PaidAt        time.Time
	IsTest        bool
	IsDeleted     bool
}
