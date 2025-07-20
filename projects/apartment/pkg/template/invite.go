package template

import (
	"bytes"
	"html/template"
)

type InviteData struct {
	Name          string
	EventName     string
	EventDate     string
	EventTime     string
	EventLocation string
	Message       string
	RSVPLink      string
	OrganizerName string
}

func NewInvite(data InviteData) ([]byte, error) {
	t, err := template.New("invite").ParseFiles("pkg/template/invite_email.html")
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "InviteEmail", data); err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}
