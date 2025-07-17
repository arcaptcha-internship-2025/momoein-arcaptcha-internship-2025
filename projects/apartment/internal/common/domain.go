package common

import (
	"regexp"

	"github.com/google/uuid"
)

type ID = uuid.UUID

var NilID = ID{}

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) IsValid() bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	r := regexp.MustCompile(emailRegex)
	return r.Match([]byte(e))
}
