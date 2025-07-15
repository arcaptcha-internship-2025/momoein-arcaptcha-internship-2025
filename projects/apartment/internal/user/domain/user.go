package domain

import (
	"errors"
	"regexp"
	"slices"

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
	Email  string
)

var NilID = UserID{}

func (e Email) IsValid() bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	r := regexp.MustCompile(emailRegex)
	return r.Match([]byte(e))
}

func (e Email) String() string {
	return string(e)
}

type User struct {
	ID        UserID
	Email     Email
	password  []byte
	FirstName string
	LastName  string
}

func New(id UserID, email, pass, firstName, lastName string) *User {
	return &User{
		ID:        id,
		Email:     Email(email),
		password:  []byte(email),
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

func (u *User) ComparePassword(pass []byte) {
	bcrypt.CompareHashAndPassword(u.password, pass)
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
	Email Email
}

func (f *UserFilter) IsValid() bool {
	return f.ID != NilID || f.Email.IsValid()
}
