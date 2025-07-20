package apartment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/template"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrApartmentOnCreate = errors.New("apartment creation failed")
	ErrInvalidAdmin      = errors.New("invalid apartment admin")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrOnInviteMember    = errors.New("error on invite member")
	ErrUnregisteredUser  = errors.New("unregistered user")
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidToken      = errors.New("invalid token")
)

type service struct {
	repo port.Repo
	mail port.Email
}

func NewService(r port.Repo, mail port.Email) port.Service {
	return &service{
		repo: r,
		mail: mail,
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
	*domain.Invite, error,
) {
	log := appctx.Logger(ctx)

	if err := s.validateApartmentAdmin(ctx, apartmentID, adminID); err != nil {
		// Unauthorized: current user is not the admin of this apartment
		if errors.Is(err, ErrInvalidAdmin) {
			log.Error("permission denied", zap.Error(err))
			return nil, fp.WrapErrors(ErrOnInviteMember, ErrPermissionDenied, err)
		}
		if errors.Is(err, ErrNotFound) {
			log.Error("apartment not found", zap.Error(err))
			return nil, fp.WrapErrors(ErrOnInviteMember, err)
		}
		log.Error("apartment admin validate failed", zap.Error(err))
		return nil, fp.WrapErrors(ErrOnInviteMember, err)
	}

	invite := &domain.Invite{
		Email:     userEmail,
		Status:    domain.InviteStatusPending,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	invite, err := s.repo.InviteMember(ctx, apartmentID, invite)
	if err != nil {
		log.Error("repo invite failed", zap.Error(err))
		return nil, fp.WrapErrors(ErrOnInviteMember, err)
	}

	// !! Dirty Code
	log.Warn("Dirty Code")

	inviteData := template.InviteData{
		Name:          invite.Email.String(),
		EventName:     "Apartment",
		Message:       "Please use the following link to accept the invitation:",
		RSVPLink:      fmt.Sprintf("http://127.0.0.1:8080/apartment/accepte/%s", invite.Token),
		OrganizerName: "The ArCaptcha Team",
	}
	msg, err := template.NewInvite(inviteData)
	if err != nil {
		log.Error("invite template", zap.Error(err))
		return nil, fp.WrapErrors(ErrOnInviteMember, err)
	}
	log.Info("mail", zap.Any("mail service", s.mail))
	err = s.mail.Send([]string{userEmail.String()}, msg)
	if err != nil {
		log.Error("send email", zap.Error(err))
		return nil, fp.WrapErrors(ErrOnInviteMember, err)
	}

	return invite, nil
}

func (s *service) validateApartmentAdmin(ctx context.Context, ApartmentID, adminID common.ID) error {
	apartment, err := s.repo.Get(ctx, &domain.ApartmentFilter{ID: ApartmentID})
	if err != nil {
		return err
	}
	if apartment == nil {
		return ErrNotFound
	}
	if apartment.AdminID != adminID {
		return ErrInvalidAdmin
	}
	return nil
}

func (s *service) AcceptInvite(ctx context.Context, token string) (*domain.Apartment, error) {
	panic("unimplemented")
}

func (s *service) Members(ctx context.Context, id common.ID) ([]domain.ApartmentMember, error) {
	panic("unimplemented")
}
