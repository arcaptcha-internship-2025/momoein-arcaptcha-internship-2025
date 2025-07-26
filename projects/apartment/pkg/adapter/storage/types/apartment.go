package types

import (
	aptDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	bilDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

type Apartment struct {
	Model
	Name       string
	Address    string
	UnitNumber int64
	AdminID    string
}

func ApartmentDomainToStorage(a *aptDomain.Apartment) *Apartment {
	return &Apartment{
		Model:      Model{ID: a.ID.String()},
		Name:       a.Name,
		Address:    a.Address,
		UnitNumber: a.UnitNumber,
		AdminID:    a.AdminID.String(),
	}
}

func ApartmentStorageToDomain(a *Apartment) *aptDomain.Apartment {
	id := common.NilID
	_ = id.UnmarshalText([]byte(a.ID))
	adminId := common.NilID
	_ = adminId.UnmarshalText([]byte(a.AdminID))
	return &aptDomain.Apartment{
		ID:         id,
		Name:       a.Name,
		Address:    a.Address,
		UnitNumber: a.UnitNumber,
		AdminID:    adminId,
		Members:    []aptDomain.ApartmentMember{},
		Bills:      []bilDomain.Bill{},
	}
}
