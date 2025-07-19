package types

import (
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

type Apartment struct {
	Model
	Name       string
	Address    string
	UnitNumber int64
	AdminID    string
}

func ApartmentDomainToStorage(a *domain.Apartment) *Apartment {
	return &Apartment{
		Model:      Model{ID: a.ID.String()},
		Name:       a.Name,
		Address:    a.Address,
		UnitNumber: a.UnitNumber,
		AdminID:    a.AdminID.String(),
	}
}

func ApartmentStorageToDomain(a *Apartment) *domain.Apartment {
	id := common.NilID
	_ = id.UnmarshalText([]byte(a.ID))
	adminId := common.NilID
	_ = adminId.UnmarshalText([]byte(a.AdminID))
	return &domain.Apartment{
		ID:         id,
		Name:       a.Name,
		Address:    a.Address,
		UnitNumber: a.UnitNumber,
		AdminID:    adminId,
		Members:    []domain.ApartmentMember{},
		Bills:      []domain.Bill{},
	}
}
