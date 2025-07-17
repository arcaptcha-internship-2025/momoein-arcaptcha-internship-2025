package appjwt

import "github.com/golang-jwt/jwt/v5"

type InviteClaims struct {
	jwt.RegisteredClaims
	IsRegisteredUser bool
}

func CreateInviteToken(secret []byte, claims *InviteClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func ParseInviteToken(tokenString string, secret []byte) (*InviteClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &InviteClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if token == nil {
		return nil, ErrNilToken
	}

	var claim *InviteClaims
	if token.Claims != nil {
		cc, ok := token.Claims.(*InviteClaims)
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
