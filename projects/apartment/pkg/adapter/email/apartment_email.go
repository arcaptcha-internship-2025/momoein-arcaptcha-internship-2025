package email

import (
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/smaila"
)

type apartmentEmail struct {
	*smaila.Sender
}

func NewApartmentEmail(endpoint string) port.EmailSender {
	return &apartmentEmail{Sender: smaila.NewSender(endpoint)}
}

func (a *apartmentEmail) Send(to []string, msg *common.EmailMessage) error {
	return a.Sender.Send(to, msg.Subject, msg.Body, msg.IsHTML)
}
