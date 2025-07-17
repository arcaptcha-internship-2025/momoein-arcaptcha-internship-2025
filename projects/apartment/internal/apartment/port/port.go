package port

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

type Service interface {
	Create(ctx context.Context, a *domain.Apartment) (*domain.Apartment, error)
	InviteMember(
		ctx context.Context,
		adminID, apartmentID common.ID,
		userEmail common.Email,
	) (
		*domain.ApartmentMember, error,
	)
	// AcceptInvite(ctx context.Context, token string) (*domain.Apartment, error)
	Members(ctx context.Context, id common.ID) ([]domain.ApartmentMember, error)
}

type Repo interface {
	Create(ctx context.Context, a *domain.Apartment) (*domain.Apartment, error)
	Get(ctx context.Context, f *domain.ApartmentFilter) (*domain.Apartment, error)
	AddMember(
		ctx context.Context,
		apartmentID common.ID,
		memberEmail common.Email,
		invite *domain.Invite,
	) (
		*domain.ApartmentMember, error,
	)
}
