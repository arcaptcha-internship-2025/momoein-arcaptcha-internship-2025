package types

import (
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

type UserBillShare struct {
	UserID       string
	BillID       string
	BillName     string
	TotalAmount  int
	MemberCount  int
	SharePerUser int
	UserPaid     int
	BalanceDue   int
}

func UserBillShareStorageToDomain(share *UserBillShare) *domain.UserBillShare {
	userID := common.NilID
	_ = userID.UnmarshalText([]byte(share.UserID))
	billID := common.NilID
	_ = billID.UnmarshalText([]byte(share.BillID))
	return &domain.UserBillShare{
		UserID:       userID,
		BillID:       billID,
		BillName:     share.BillName,
		TotalAmount:  share.TotalAmount,
		MemberCount:  share.MemberCount,
		SharePerUser: share.SharePerUser,
		UserPaid:     share.UserPaid,
		BalanceDue:   share.BalanceDue,
	}
}

func UserBillShareDomainToStorage(share *domain.UserBillShare) *UserBillShare {
	return &UserBillShare{
		UserID:       share.UserID.String(),
		BillID:       share.BillID.String(),
		BillName:     share.BillName,
		TotalAmount:  share.TotalAmount,
		MemberCount:  share.MemberCount,
		SharePerUser: share.SharePerUser,
		UserPaid:     share.UserPaid,
		BalanceDue:   share.BalanceDue,
	}
}
