package invoice

import (
	"context"
	"time"

	invoiceentity "github.com/albertwidi/go_project_example/entity/invoice"
	"github.com/albertwidi/go_project_example/internal/pkg/sqldb"
	"github.com/internal/pkg/pq"
)

// Repository of invoice
type Repository struct {
	db *sqldb.DB
}

// Invoice struct
type Invoice struct {
	// ID is PRIMARY_KEY
	ID string `db:"id"`
	// Number/invoice_number might be indexed
	Number        string `db:"invoice_number"`
	OrderID       string `db:"order_id"`
	InvoiceFrom   int64  `db:"invoice_from"`
	InvoiceTO     int64  `db:"invoice_to"`
	Type          int    `db:"invoice_type"`
	Total         int64  `db:"total"`
	DiscountTotal int64  `db:"discount_total"`
	// grand total = total - discount_total
	GrandTotal  int64                `db:"grand_total"`
	Details     []string             `db:"invoice_details"`
	Status      invoiceentity.Status `db:"invoice_status"`
	Description string               `db:"description"`
	DueDate     time.Time            `db:"due_date"`
	PaidAt      pq.NullTime          `db:"paid_at"`
	CreatedAt   time.Time            `db:"created_at"`
	UpdatedAt   pq.NullTime          `db:"updated_at"`
	CreatedBy   int64                `db:"created_by"`
	IsTest      bool                 `db:"is_test"`
	IsDeleted   bool                 `db:"is_deleted"`
}

// Detail for invoice detail
type Detail struct {
	// ID is PRIMARY_KEY
	ID string `db:"id"`
	// InvoiceID is INDEXED
	InvoiceID string `db:"invoice_id"`
	Amount    int64  `db:"amount"`
	Discount  int64  `db:"discount"`
	// total = amount - discount
	Total        int64     `db:"total"`
	ItemName     string    `db:"item_name"`
	ItemQuantity int64     `db:"item_quantity"`
	Description  string    `db:"description"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	IsTest       bool      `db:"is_test"`
	IsDeleted    bool      `db:"is_deleted"`
}

// PaidInvoice struct
// store all data of paid invoice
type PaidInvoice struct {
	// InvoiceID is PRIMARY_KEY
	InvoiceID     string    `db:"invoice_id"`
	InvoiceNumber string    `db:"invoice_number"`
	PaymentID     string    `db:"payment_id"`
	OrderID       string    `db:"order_id"`
	InvoiceFrom   int64     `db:"invoice_from"`
	InvoiceTo     int64     `db:"invoice_to"`
	CreatedAt     time.Time `db:"created_at"`
	CreatedBy     int64     `db:"created_by"`
	PaidAt        time.Time `db:"paid_at"`
	IsTest        bool      `db:"is_test"`
	IsDeleted     bool      `db:"is_deleted"`
}

// New invoice repo
func New(db *sqldb.DB) *Repository {
	r := Repository{
		db: db,
	}
	return &r
}

const (
	createNewInvoiceQuery = `
INSERT INTO invoices (
	id,
	invoice_number,
	order_id,
	invoice_from,
	invoice_to,
	invoice_type,
	total,
	discount_total,
	grand_total,
	invoice_details,
	invoice_status,
	description,
	due_date,
	created_at,
	created_by,
	is_test,
	is_deleted
) VALUES(
	?, # id
	?, # invoice_number
	?, # order_id
	?, # invoice_from
	?, # invoice_to
	?, # invoice_type
	?, # total
	?, # discount_total
	?, # grand_total
	?, # invoice_details
	?, # invoice_status
	?, # description
	?, # due_date
	?, # created_at
	?, # created_by
	?, # is_test
	?  # is deleted
)
`

	createNewInvoiceDetailQuery = `
INSERT INTO invoices_details (
	invoice_id,
	amount,
	discount,
	item_name,
	description,
	created_at,
	is_test,
	is_deleted
) VALUES(
	?, # invoice_id
	?, # amount
	?, # discount
	?, # item_name
	?, # description
	?, # created_at
	?, # is_test
	?  # is_deleted
)
`
)

// Create a new invoice
func (r Repository) Create(ctx context.Context, invoice invoiceentity.Invoice) error {
	return nil
}

const (
	createNewPaidInvoice = `
INSERT INTO paid_invoices (
	invoice_id,
	invoice_number,
	payment_id,
	order_id,
	invoice_from,
	invoice_to,
	created_at,
	created_by
	paid_at,
	is_test,
	is_deleted
) VALUES (
	?, # invoice_id
	?, # invoice_number
	?, # payment_id
	?, # order_id
	?, # invoice_from
	?, # invoice_to
	?, # created_at
	?, # paid_at
	?, # is_test
	?  # is_deleted
)
`
)

// CreatePaid to create a new paid invoice
func (r Repository) CreatePaid(ctx context.Context, paidInvoice invoiceentity.PaidInvoice) error {
	return nil
}

// Get invoice
func (r Repository) Get(ctx context.Context, invoiceID string, withDetails bool) (invoiceentity.Invoice, error) {
	i := invoiceentity.Invoice{}
	return i, nil
}

// UpdateStatus to update the status of invoice
func (r Repository) UpdateStatus(ctx context.Context, invoiceID string, status invoiceentity.Status) error {
	return nil
}

// UpdatePaidAt to paid
func (r Repository) UpdatePaidAt(ctx context.Context, invoiceID string, paidAt time.Time) error {
	return nil
}
