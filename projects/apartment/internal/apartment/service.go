package apartment

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
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
	ErrOnAcceptInvite    = errors.New("error on accept invite")
	ErrExpiredToken      = errors.New("expired token")
	ErrOnSendEmail       = errors.New("failed to send email")
	ErrOnParsURL         = errors.New("failed to parse url")
	ErrOnGenerateMessage = errors.New("failed to generate message")
	ErrInvalidEmail      = errors.New("invalid email")
)

type service struct {
	repo port.Repo
	mail port.EmailSender
}

func NewService(r port.Repo, mail port.EmailSender) port.Service {
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
	acceptURL string,
) (
	*domain.Invite, error,
) {
	log := appctx.Logger(ctx)

	if !userEmail.IsValid() {
		log.Error("invalid email")
		return nil, fp.WrapErrors(ErrOnInviteMember, ErrInvalidEmail)
	}

	if err := s.validateApartmentAdmin(ctx, apartmentID, adminID); err != nil {
		log.Error("apartment admin validation failed", zap.Error(err))
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

	to := []string{invite.Email.String()}
	msg, err := s.generateInviteMessage(invite, acceptURL)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnInviteMember, err)
	}

	err = s.mail.Send(to, msg)
	if err != nil {
		log.Error("send email", zap.Error(err))
		return nil, fp.WrapErrors(ErrOnInviteMember, ErrOnSendEmail, err)
	}

	return invite, nil
}

func (s *service) validateApartmentAdmin(
	ctx context.Context, ApartmentID, adminID common.ID,
) error {
	apartment, err := s.repo.Get(ctx, &domain.ApartmentFilter{ID: ApartmentID})
	if err != nil {
		return err
	}
	if apartment == nil {
		return ErrNotFound
	}
	if apartment.AdminID != adminID {
		return fp.WrapErrors(ErrPermissionDenied, ErrInvalidAdmin)
	}
	return nil
}

func (s *service) generateInviteMessage(
	invite *domain.Invite, acceptURL string,
) (
	*common.EmailMessage, error,
) {
	rsvpLink, err := url.Parse(acceptURL)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnParsURL, err)
	}
	rsvpLink.RawQuery = fmt.Sprintf("%s=%s", "token", invite.Token)

	inviteData := template.InviteData{
		Name:          strings.Split(invite.Email.String(), "@")[0],
		ApartmentName: "Apartment",
		Message:       "Please use the following link to accept the invitation:",
		RSVPLink:      rsvpLink.String(),
		OrganizerName: "The ArCaptcha Team",
	}
	body, err := template.NewInvite(inviteData)
	if err != nil {
		return nil, err
	}

	return &common.EmailMessage{
		Subject: "apartment invite",
		Body:    body,
		IsHTML:  true,
	}, nil
}

func (s *service) AcceptInvite(ctx context.Context, token string) error {
	if err := common.ValidateID(token); err != nil {
		return fp.WrapErrors(ErrOnAcceptInvite, ErrInvalidToken, err)
	}
	err := s.repo.AcceptInvite(ctx, token)
	if err != nil {
		return fp.WrapErrors(ErrOnAcceptInvite, err)
	}
	return nil
}

func (s *service) Members(ctx context.Context, id common.ID) ([]domain.ApartmentMember, error) {
	panic("unimplemented")
}
