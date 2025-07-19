package appjwt

import (
	"errors"

	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNilToken     = errors.New("invalid token (nil)")
	ErrInvalidToken = errors.New("token is not valid")
)

const (
	UserEmailKey appctx.CtxKey = "UserEmail"
	UserIDKey    appctx.CtxKey = "UserID"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    string
	UserEMail string
}

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if token == nil {
		return nil, ErrNilToken
	}

	var claim *UserClaims
	if token.Claims != nil {
		cc, ok := token.Claims.(*UserClaims)
		if ok {
			claim = cc
		}
	}

	if err != nil {
		return claim, err
	}

	if !token.Valid {
		return claim, ErrInvalidToken
	}

	return claim, nil
}
