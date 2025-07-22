package domain

import (
	"errors"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

var (
	ErrBillMissingType       = errors.New("bill type is required")
	ErrBillInvalidBillNumber = errors.New("bill number must be greater than zero")
	ErrBillNegativeAmount    = errors.New("amount cannot be negative")
	ErrBillMissingDueDate    = errors.New("due date is required")
)

type BillType string

func (bt BillType) String() string {
	return string(bt)
}

const (
	BillElectricity BillType = "electricity"
	BillWater       BillType = "water"
	BillGas         BillType = "gas"
)

var validBillTypes = map[BillType]struct{}{
	BillElectricity: {},
	BillWater:       {},
	BillGas:         {},
}

func (bt BillType) IsValid() bool {
	_, ok := validBillTypes[bt]
	return ok
}

type Image struct {
	Name    string
	Type    string
	Size    int64
	Content []byte
}

type Bill struct {
	ID          common.ID
	Name        string
	Type        BillType
	BillNumber  int64
	DueDate     time.Time
	Amount      int64
	Status      PaymentStatus
	PaidAt      time.Time
	Image       Image
	HasImage    bool
	ImageID     common.ID
	ApartmentID common.ID
}

func (b *Bill) Validate() error {
	if b.Type == "" {
		return ErrBillMissingType
	}
	if b.BillNumber <= 0 {
		return ErrBillInvalidBillNumber
	}
	if b.Amount < 0 {
		return ErrBillNegativeAmount
	}
	if b.DueDate.IsZero() {
		return ErrBillMissingDueDate
	}
	return nil
}

type BillFilter struct {
	ID          common.ID
	ApartmentID common.ID
	Type        BillType
	BillNumber  int64
}

func (b Bill) IsPaid() bool {
	return b.Status == PaymentStatusPaid && !b.PaidAt.IsZero()
}

type PaymentStatus string

const (
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusOverdue PaymentStatus = "overdue"
)

func (ps PaymentStatus) String() string {
	return string(ps)
}

var validPaymentStatuses = map[PaymentStatus]struct{}{
	PaymentStatusUnpaid:  {},
	PaymentStatusPaid:    {},
	PaymentStatusOverdue: {},
}

func (ps PaymentStatus) IsValid() bool {
	_, ok := validPaymentStatuses[ps]
	return ok
}

type Payment struct {
	ID      common.ID
	BillID  common.ID
	PayerID common.ID
	Amount  int64
	PaidAt  time.Time
}
