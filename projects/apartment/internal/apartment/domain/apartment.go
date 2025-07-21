package domain

import (
	"time"

	billDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
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
	ID        common.ID
	Email     common.Email
	Status    InviteStatus
	Token     string
	ExpiresAt time.Time
}

type ApartmentMember struct {
	userDomain.User
	Debt int64
}

type Apartment struct {
	ID         common.ID
	Name       string
	Address    string
	UnitNumber int64
	AdminID    common.ID
	Members    []ApartmentMember
	Bills      []billDomain.Bill
}

func (a *Apartment) Validate() error {
	return nil
}

type ApartmentFilter struct {
	ID common.ID
}
