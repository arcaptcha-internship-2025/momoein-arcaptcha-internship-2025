package types

import (
	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
)

type User struct {
	ID        string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func UserDomainToStorage(u *userDomain.User) *User {
	return &User{
		ID:        u.ID.String(),
		Email:     u.Email.String(),
		Password:  string(u.Password()),
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func UserStorageToDomain(u *User) *userDomain.User {
	user := &userDomain.User{
		ID:        userDomain.UserID([]byte(u.ID)),
		Email:     userDomain.Email(u.Email),
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
	user.SetPassword([]byte(u.Password))
	return user
}
