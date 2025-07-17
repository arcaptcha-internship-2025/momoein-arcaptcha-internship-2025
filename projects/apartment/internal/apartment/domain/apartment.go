package domain

import (
	"image"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
)

type InviteStatus string

func (s InviteStatus) String() string {
	return string(s)
}

const (
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusDeclined InviteStatus = "declined"
	InviteStatusExpired  InviteStatus = "expired"
)

var validInviteStatus = map[InviteStatus]struct{}{
	InviteStatusAccepted: {},
	InviteStatusPending:  {},
	InviteStatusDeclined: {},
	InviteStatusExpired:  {},
}

func (s InviteStatus) IsValid() (ok bool) {
	_, ok = validInviteStatus[s]
	return
}

type Invite struct {
	Status    InviteStatus
	Token     string
	ExpiresAt time.Time
}

type ApartmentMember struct {
	domain.User
	Invite Invite
	Debt   int64
}

type Apartment struct {
	ID         common.ID
	Name       string
	Address    string
	UnitNumber int64
	AdminID    common.ID
	Members    []ApartmentMember
	Bills      []Bill
}

func (a *Apartment) Validate() error {
	return nil
}

type ApartmentFilter struct {
	ID common.ID
}



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

type Bill struct {
	ID          common.ID
	Name        string
	Type        BillType
	BillNumber  int64
	DueDate     time.Time
	Amount      int64
	Status      PaymentStatus
	PaidAt      time.Time
	Image       image.Image
	ImageID     common.ID
	ApartmentID common.ID
}

func (b Bill) IsPaid() bool {
	return b.Status == PaymentStatusPaid && !b.PaidAt.IsZero()
}

type Payment struct {
	ID      common.ID
	BillID  common.ID
	PayerID common.ID
	Amount  int64
	PaidAt  time.Time
}
