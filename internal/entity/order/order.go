package order

import "time"

// Status of order
type Status int

// Type of order
type Type int

// Order struct
type Order struct {
	ID            string
	UserID        int64
	OrderType     int32
	Total         int64
	DiscountTotal int64
	RefundTotal   int64
	CashbackTotal int64
	// grand total = total - discount total - refund total
	GrandTotal       int64
	Status           Status
	Details          []Detail
	IsRefundable     bool
	RefundableBefore time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UpdatedBy        int64
	IsTest           bool
	IsDeleted        bool
}

// Detail of order
type Detail struct {
	ID             string
	OrderID        string
	ItemID         string
	ItemName       string
	ItemPrice      int64
	ItemQuantity   int64
	ItemPriceTotal int64
	DiscountAmount int64
	RefundAmount   int64
	CashbackAmount int64
	// total = item price total - discount amount - refund amount
	Total     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy int64
	IsTest    bool
	IsDeleted bool
}

// Coupons of order, to track the usage of coupons
type Coupons struct {
	OrderID        string
	UserID         int64
	CouponCode     string
	DiscountAmount int64
	// to track whether the coupon is really applied or not
	IsApplied bool
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy int64
	IsTest    bool
	IsDeleted bool
}

// RefundStatus of order
type RefundStatus int

// Refund of order
type Refund struct {
	ID          string
	OrderID     string
	Status      RefundStatus
	TotalAmount int64
	Details     []RefundDetail
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UpdatedBy   int64
	IsTest      bool
	IsDeleted   bool
}

// RefundDetail struct
type RefundDetail struct {
	RefundID      string
	OrderDetailID string
	Amount        int64
	CreatedAt     time.Time
	IsTest        bool
	IsDeleted     bool
}
