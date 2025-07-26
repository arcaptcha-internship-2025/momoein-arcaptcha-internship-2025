package template

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed invite_email.html
var inviteEmailTemplate string

type InviteData struct {
	Name          string
	EventName     string
	Message       string
	RSVPLink      string
	OrganizerName string
}

func NewInvite(data InviteData) ([]byte, error) {
	tmpl, err := template.New("InviteEmail").Parse(inviteEmailTemplate)
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}
