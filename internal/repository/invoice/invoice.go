package invoice

import (
	"context"
	"time"

	invoiceentity "github.com/albertwidi/go-project-example/internal/entity/invoice"
	"github.com/albertwidi/go-project-example/internal/pkg/sqldb"
	"github.com/lib/pq"
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
	InvoiceFrom   string `db:"invoice_from"`
	InvoiceTo     string `db:"invoice_to"`
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
	ItemName     string    `db:"item_name"`
	ItemQuantity int64     `db:"item_quantity"`
	Description  string    `db:"description"`
	CreatedAt    time.Time `db:"created_at"`
	CreatedBy    string    `db:"created_by"`
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
	:invoice_number # invoice_number
	:order_id # order_id
	:invoice_from, # invoice_from
	:invoice_to, # invoice_to
	:invoice_type, # invoice_type
	:total, # total
	:discount_total, # discount_total
	:grand_total, # grand_total
	:invoice_details, # invoice_details
	:invoice_status, # invoice_status
	:description, # description
	:due_date, # due_date
	:created_at, # created_at
	:created_by, # created_by
	:is_test, # is_test
	:is_deleted  # is deleted
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
	:invoice_id, # invoice_id
	:amount, # amount
	:discount, # discount
	:item_name, # item_name
	:description, # description
	:created_at, # created_at
	:is_test, # is_test
	:is_deleted  # is_deleted
)
`
)

// Create a new invoice
func (r *Repository) Create(ctx context.Context, invoice *invoiceentity.Invoice) error {
	var details []string

	// create database transaction for creating invoice and invoice detail
	tx, err := r.db.Leader().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// create the detail first
	for idx, invD := range invoice.Details {
		invoiceDetail := Detail{
			InvoiceID:    invD.InvoiceID,
			Amount:       invD.Amount,
			Discount:     invD.Discount,
			ItemName:     invD.ItemName,
			ItemQuantity: invD.ItemQuantity,
			Description:  invD.Description,
			CreatedAt:    time.Now(),
			CreatedBy:    invD.CreatedBy,
			IsTest:       invD.IsTest,
		}
		qDetail, args, err := r.db.BindNamed(createNewInvoiceDetailQuery, invoiceDetail)
		if err != nil {
			return err
		}

		rows, err := tx.QueryContext(ctx, qDetail, args)
		if err != nil {
			return err
		}
		// only return 1 result: ID
		if rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return err
			}
			details = append(details, id)
			// set invoice details id
			invoice.Details[idx].ID = id
		}
	}

	// insert the invoice
	inv := Invoice{
		Number:        invoice.Number,
		OrderID:       invoice.OrderID,
		InvoiceFrom:   invoice.InvoiceFrom,
		InvoiceTo:     invoice.InvoiceTo,
		Type:          invoice.Type,
		Total:         invoice.Total,
		DiscountTotal: invoice.DiscountTotal,
		GrandTotal:    invoice.GrandTotal,
		Details:       details,
		Status:        invoice.Status,
		Description:   invoice.Description,
		DueDate:       invoice.DueDate,
		CreatedAt:     time.Now(),
		CreatedBy:     invoice.CreatedBy,
		IsTest:        invoice.IsTest,
	}
	qInv, args, err := r.db.BindNamed(createNewInvoiceQuery, inv)
	if err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, qInv, args)
	if err != nil {
		return err
	}
	if rows.Next() {
		if err := rows.Scan(invoice.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
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
