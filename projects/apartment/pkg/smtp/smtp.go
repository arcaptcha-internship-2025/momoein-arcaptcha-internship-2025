package smtp

import (
	"fmt"
	"net/smtp"
)

type smtpService struct {
	host     string
	port     uint
	from     string
	username string
	password string
}

func NewSMTPService(host string, port uint, from, username, password string) *smtpService {
	return &smtpService{
		host:     host,
		port:     port,
		from:     from,
		username: username,
		password: password,
	}
}

func (s *smtpService) Send(to []string, msg []byte) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, to, msg)
}
