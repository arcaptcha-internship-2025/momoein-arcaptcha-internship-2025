package storage

import (
	"context"
	"database/sql"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/storage/types"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"go.uber.org/zap"
)

type apartmentRepo struct {
	db *sql.DB
}

func NewApartmentRepo(db *sql.DB) port.Repo {
	return &apartmentRepo{db: db}
}

func (r *apartmentRepo) Create(ctx context.Context, a *domain.Apartment) (*domain.Apartment, error) {
	log := appctx.Logger(ctx)

	ap := types.ApartmentDomainToStorage(a)

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO apartments(name, address, unit_number, admin_id)
		VALUES($1, $2, $3, $4)
		RETURNING id;`,
		ap.Name, ap.Address, ap.UnitNumber, ap.AdminID,
	).Scan(&ap.ID)

	if err != nil {
		log.Error("failed to execute query", zap.Error(err))
		return nil, err
	}

	return types.ApartmentStorageToDomain(ap), nil
}

func (r *apartmentRepo) Get(ctx context.Context, f *domain.ApartmentFilter) (*domain.Apartment, error) {
	panic("unimplemented")
}

func (r *apartmentRepo) AddMember(
	ctx context.Context,
	apartmentID common.ID,
	memberEmail common.Email,
	invite *domain.Invite,
) (
	*domain.ApartmentMember, error,
) {
	panic("unimplemented")
}
