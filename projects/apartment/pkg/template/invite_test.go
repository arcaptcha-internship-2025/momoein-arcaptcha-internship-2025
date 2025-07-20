package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInvite(t *testing.T) {
	data := InviteData{
		Name:          "Alex",
		EventName:     "Go Developers Meetup 2025",
		EventDate:     "August 12, 2025",
		EventTime:     "6:00 PM",
		EventLocation: "Tech Hub, Downtown",
		Message:       "Join us for an evening of networking, talks, and pizza!",
		RSVPLink:      "https://example.com/rsvp",
		OrganizerName: "The Go Team",
	}

	msg, err := NewInvite(data)
	assert.NoError(t, err)
	_ = msg
}
