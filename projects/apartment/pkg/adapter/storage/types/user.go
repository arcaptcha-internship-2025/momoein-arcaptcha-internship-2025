package types

import (
	"database/sql"

	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	"github.com/google/uuid"
)

type User struct {
	ID        string
	Email     string
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
}

func UserDomainToStorage(u *userDomain.User) *User {
	firstName := sql.NullString{
		String: u.FirstName,
		Valid:  len(u.FirstName) > 0,
	}
	lastName := sql.NullString{
		String: u.LastName,
		Valid:  len(u.LastName) > 0,
	}
	return &User{
		ID:        u.ID.String(),
		Email:     u.Email.String(),
		Password:  string(u.Password()),
		FirstName: firstName,
		LastName:  lastName,
	}
}

func UserStorageToDomain(u *User) *userDomain.User {
	var id = uuid.Nil
	if err := uuid.Validate(u.ID); err == nil {
		id = uuid.MustParse(u.ID)
	}
	return userDomain.New(
		id, u.Email, u.Password,
		u.FirstName.String, u.LastName.String,
	)
}
