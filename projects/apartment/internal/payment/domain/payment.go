package domain

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the current status of a payment
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentPaid      PaymentStatus = "paid"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCancelled PaymentStatus = "cancelled"
)

type CallbackData map[string]any

type Payment struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	BillID      uuid.UUID `json:"bill_id"`
	PayerID     uuid.UUID `json:"payer_id"`
	Amount      int64     `json:"amount"` // in smallest currency unit, e.g., Rial or Toman
	PaymentDate time.Time `json:"payment_date"`

	Status        PaymentStatus `json:"status"`
	Gateway       string        `json:"gateway"`
	TransactionID string        `json:"transaction_id,omitempty"`
	CallbackData  CallbackData  `json:"callback_data,omitempty"` // Use map for parsed JSONB
}
