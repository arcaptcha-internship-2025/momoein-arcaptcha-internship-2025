package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
)

var (
	ErrUserOnValidate = errors.New("user validation failed")
	ErrUserOnCreate   = errors.New("user Creation failed")
)

type service struct {
	repo port.Repo
}

func NewService(r port.Repo) port.Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if err := u.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUserOnValidate, err)
	}
	user, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUserOnCreate, err)
	}
	return user, nil
}

func (s *service) Get(context.Context, *domain.UserFilter) (*domain.User, error) {
	panic("unimplemented")
}

func (s *service) Delete(context.Context, *domain.UserFilter) error {
	panic("unimplemented")
}
