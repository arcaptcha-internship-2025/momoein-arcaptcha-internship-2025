package domain

import (
	"errors"
	"slices"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrUserShortPassword = errors.New("password must be at least 8 characters long")
	ErrUserLongPassword  = errors.New("password must be at most 72 characters long")
)

type (
	UserID = uuid.UUID
)

var NilID = UserID{}

type User struct {
	ID        UserID
	Email     common.Email
	password  []byte
	FirstName string
	LastName  string
}

func NewUser(id UserID, email, pass, firstName, lastName string) *User {
	return &User{
		ID:        id,
		Email:     common.Email(email),
		password:  []byte(pass),
		FirstName: firstName,
		LastName:  lastName,
	}
}

func (u *User) Password() []byte {
	return slices.Clone(u.password)
}

func (u *User) SetPassword(pass []byte) error {
	p, err := bcrypt.GenerateFromPassword(pass, 12)
	if err != nil {
		return err
	}
	u.password = p
	return nil
}

func (u *User) ComparePassword(pass []byte) error {
	return bcrypt.CompareHashAndPassword(u.password, pass)
}

func (u *User) Validate() error {
	switch {
	case !u.Email.IsValid():
		return ErrInvalidEmail
	case len(u.password) < 8:
		return ErrUserShortPassword
	case len(u.password) > 72:
		return ErrUserLongPassword
	default:
		return nil
	}
}

type UserFilter struct {
	ID    UserID
	Email common.Email
}

func (f *UserFilter) IsValid() bool {
	return f.ID != NilID || f.Email.IsValid()
}
