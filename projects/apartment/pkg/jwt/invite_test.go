package appjwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInviteToken(t *testing.T) {
	secret := []byte("secret")
	ic := &InviteClaims{IsRegisteredUser: false}

	token, err := CreateInviteToken(secret, ic)
	assert.NoError(t, err)

	fmt.Printf("len: %d, token: %q\n", len(token), token)

	claim, err := ParseInviteToken(token, secret)
	assert.NoError(t, err)
	assert.IsType(t, ic, claim)
	assert.Equal(t, ic.IsRegisteredUser, claim.IsRegisteredUser)
}
