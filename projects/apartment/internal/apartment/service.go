package apartment

import (
	"context"
	"errors"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
	"github.com/google/uuid"
)

var (
	ErrApartmentOnCreate = errors.New("apartment creation failed")
	ErrInvalidAdmin      = errors.New("invalid apartment admin")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrOnInviteMember    = errors.New("error on invite member")
	ErrUnregisteredUser  = errors.New("unregistered user")
)

type service struct {
	repo port.Repo
}

func NewService(r port.Repo) port.Service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, a *domain.Apartment) (*domain.Apartment, error) {
	if err := a.Validate(); err != nil {
		return nil, fp.WrapErrors(ErrApartmentOnCreate, err)
	}
	apartment, err := s.repo.Create(ctx, a)
	if err != nil {
		return nil, fp.WrapErrors(ErrApartmentOnCreate, err)
	}
	return apartment, nil
}

func (s *service) InviteMember(
	ctx context.Context,
	adminID, apartmentID common.ID,
	userEmail common.Email,
) (
	*domain.ApartmentMember, error,
) {
	if err := s.validateApartmentAdmin(ctx, apartmentID, adminID); err != nil {
		// Unauthorized: current user is not the admin of this apartment
		if errors.Is(err, ErrInvalidAdmin) {
			return nil, fp.WrapErrors(ErrOnInviteMember, ErrPermissionDenied, err)
		}
		return nil, fp.WrapErrors(ErrOnInviteMember, err)
	}

	invite := &domain.Invite{
		Status:    domain.InviteStatusPending,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	member, err := s.repo.AddMember(ctx, apartmentID, userEmail, invite)
	unregistered := errors.Is(err, ErrUnregisteredUser) || member.ID == common.NilID
	if err != nil && !unregistered {
		return nil, fp.WrapErrors(ErrOnInviteMember, err)
	}
	if unregistered {
		// log
	}

	// TODO: send invitation
	return member, nil
}

func (s *service) validateApartmentAdmin(ctx context.Context, ApartmentID, adminID common.ID) error {
	apartment, err := s.repo.Get(ctx, &domain.ApartmentFilter{ID: ApartmentID})
	if err != nil {
		return err
	}
	if apartment.AdminID != adminID {
		return ErrInvalidAdmin
	}
	return nil
}

func (s *service) Members(ctx context.Context, id common.ID) ([]domain.ApartmentMember, error) {
	panic("unimplemented")
}
