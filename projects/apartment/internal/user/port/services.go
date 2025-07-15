package port

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
)

type Service interface {
	Create(context.Context, *domain.User) (*domain.User, error)
	Get(context.Context, *domain.UserFilter) (*domain.User, error)
	Delete(context.Context, *domain.UserFilter) error
}
