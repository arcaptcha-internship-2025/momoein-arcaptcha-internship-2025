package types

import (
	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	"github.com/google/uuid"
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
	var id = uuid.Nil
	if err := uuid.Validate(u.ID); err == nil {
		id = uuid.MustParse(u.ID)
	}
	return userDomain.New(id, u.Email, u.Password, u.FirstName, u.LastName)
}
