package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
)

var (
	ErrUserOnValidate     = errors.New("user validation failed")
	ErrUserOnCreate       = errors.New("user Creation failed")
	ErrInvalidOrNilFilter = errors.New("invalid or empty filter")
	ErrUserOnGet          = errors.New("user retrieve failed")
	ErrUserOnDelete       = errors.New("error on deleting failed")
	ErrUserNotFound       = errors.New("user not found")
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

func (s *service) Get(ctx context.Context, filter *domain.UserFilter) (*domain.User, error) {
	if !filter.IsValid() {
		return nil, ErrInvalidOrNilFilter
	}
	u, err := s.repo.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUserOnGet, err)
	}
	return u, nil
}

func (s *service) Delete(ctx context.Context, filter *domain.UserFilter) error {
	if !filter.IsValid() {
		return ErrInvalidOrNilFilter
	}
	err := s.repo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserOnDelete, err)
	}
	return nil
}
