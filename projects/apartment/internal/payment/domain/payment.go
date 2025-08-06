package domain

import (
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

// PaymentStatus represents the current status of a payment
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentPaid      PaymentStatus = "paid"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCancelled PaymentStatus = "cancelled"
)

func (ps PaymentStatus) String() string {
	return string(ps)
}

var validPaymentStatuses = map[PaymentStatus]struct{}{
	PaymentPending:   {},
	PaymentPaid:      {},
	PaymentFailed:    {},
	PaymentCancelled: {},
}

func (ps PaymentStatus) IsValid() bool {
	_, ok := validPaymentStatuses[ps]
	return ok
}

type CallbackData map[string]any

type Payment struct {
	ID        common.ID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`

	BillID  common.ID `json:"billId"`
	PayerID common.ID `json:"payerId"`
	Amount  int64     `json:"amount"`
	PaidAt  time.Time `json:"paymentDate"`

	Status        PaymentStatus `json:"status"`
	Gateway       string        `json:"gateway"`
	TransactionID string        `json:"transactionId,omitempty"`
	CallbackData  CallbackData  `json:"callbackData,omitempty"` // Use map for parsed JSONB
}

type BillWithAmount struct {
	BillID common.ID
	Amount int64
}

type GatewayType string

func (g GatewayType) String() string {
	return string(g)
}

const MockGateway = "mock-gateway"

var validGateways = map[GatewayType]struct{}{
	MockGateway: {},
}

func (g GatewayType) IsValid() bool {
	_, ok := validGateways[g]
	return ok
}

type Transaction struct {
	PaymentIDs  []common.ID
	Amount      int64             // Amount to be paid
	PayerID     common.ID         // Who is paying
	Bills       []BillWithAmount  // which bill this is for
	CallbackURL string            // URL to hit after payment (e.g. /payment/callback)
	Metadata    map[string]string // Optional: custom key-value data for tracking
}
