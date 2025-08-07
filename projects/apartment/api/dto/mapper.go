package dto

import (
	apartmentDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	billDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	paymentd "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
)

func UserDTOToDomain(u *User) *userDomain.User {
	id := userDomain.NilID
	id.UnmarshalText([]byte(u.ID))
	return userDomain.NewUser(id, u.Email, u.Password, u.FirstName, u.LastName)
}

func ApartmentDTOToDomain(a *Apartment) *apartmentDomain.Apartment {
	id := common.NilID
	_ = id.UnmarshalText([]byte(a.ID))
	adminId := common.NilID
	_ = adminId.UnmarshalText([]byte(a.AdminID))
	return &apartmentDomain.Apartment{
		ID:         id,
		Name:       a.Name,
		Address:    a.Address,
		UnitNumber: a.UnitNumber,
		AdminID:    adminId,
		Members:    []apartmentDomain.ApartmentMember{},
		Bills:      []billDomain.Bill{},
	}
}

func ApartmentDomainToDTO(a *apartmentDomain.Apartment) *Apartment {
	return &Apartment{
		ID:         a.ID.String(),
		Name:       a.Name,
		Address:    a.Address,
		UnitNumber: a.UnitNumber,
		AdminID:    a.AdminID.String(),
	}
}

func RedirectGatewayDomainToDTO(rg *paymentd.RedirectGateway) *RedirectGateway {
	return &RedirectGateway{
		Method: rg.Method,
		URL:    rg.URL,
		Body:   rg.Body,
	}
}
