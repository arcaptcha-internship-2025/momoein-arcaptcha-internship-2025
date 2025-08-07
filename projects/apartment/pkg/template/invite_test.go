package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInvite(t *testing.T) {
	data := InviteData{
		Name:          "Alex",
		ApartmentName: "Go Developers Meetup 2025",
		Message:       "Join us for an evening of networking, talks, and pizza!",
		RSVPLink:      "https://example.com/rsvp",
		OrganizerName: "The Go Team",
	}

	msg, err := NewInvite(data)
	assert.NoError(t, err)
	_ = msg
}
