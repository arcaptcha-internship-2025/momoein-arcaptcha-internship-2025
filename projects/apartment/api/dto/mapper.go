package dto

import userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"

func UserDTOToDomain(u *User) *userDomain.User {
	id := userDomain.NilID
	id.UnmarshalText([]byte(u.ID))
	return userDomain.NewUser(id, u.Email, u.Password, u.FirstName, u.LastName)
}