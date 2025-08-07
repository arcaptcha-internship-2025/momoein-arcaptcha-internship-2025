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
	ErrMissingApartmentID    = errors.New("apartment id is required")
	ErrMissingPaymentStatus  = errors.New("payment status is required")
	ErrInvalidPaymentStatus  = errors.New("invalid payment status")
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
	Path    string
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
	Image       *Image
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
	if b.ApartmentID == common.NilID {
		return ErrMissingApartmentID
	}
	if len(b.Status) == 0 {
		return ErrMissingPaymentStatus
	}
	if !b.Status.IsValid() {
		return ErrInvalidPaymentStatus
	}
	return nil
}

func (b *Bill) SetName(s string) *Bill {
	b.Name = s
	return b
}

func (b *Bill) SetType(t BillType) *Bill {
	b.Type = t
	return b
}

func (b *Bill) SetBillNumber(n int64) *Bill {
	b.BillNumber = n
	return b
}

func (b *Bill) SetDueDate(t time.Time) *Bill {
	b.DueDate = t
	return b
}

func (b *Bill) SetAmount(i int64) *Bill {
	b.Amount = i
	return b
}

func (b *Bill) SetStatus(s PaymentStatus) *Bill {
	b.Status = s
	return b
}

func (b *Bill) SetPaidAt(t time.Time) *Bill {
	b.PaidAt = t
	return b
}

// SetImage sets the given *Image and sets the HasImage flag to true if img is not nil.
func (b *Bill) SetImage(i *Image) *Bill {
	b.Image = i
	b.HasImage = i != nil
	return b
}

func (b *Bill) SetHasImage(hi bool) *Bill {
	b.HasImage = hi
	return b
}

func (b *Bill) SetImageID(id common.ID) *Bill {
	b.ImageID = id
	return b
}

func (b *Bill) SetApartmentID(id common.ID) *Bill {
	b.ApartmentID = id
	return b
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

type UserBillShare struct {
	UserID       common.ID `json:"userId"`
	BillID       common.ID `json:"billId"`
	BillName     string    `json:"billName"`
	TotalAmount  int       `json:"totalAmount"`
	MemberCount  int       `json:"memberCount"`
	SharePerUser int       `json:"sharePerUser"`
	UserPaid     int       `json:"userPaid"`
	BalanceDue   int       `json:"balanceDue"`
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
