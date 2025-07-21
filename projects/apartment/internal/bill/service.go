package bill

import (
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
)

type service struct {
	repo port.Repo
}

func NewService(r port.Repo) port.Service {
	return &service{repo: r}
}
