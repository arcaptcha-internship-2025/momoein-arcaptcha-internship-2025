package storage

import (
	"context"
	"database/sql"
	"errors"

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
	query := `
		SELECT id, created_at, updated_at, deleted_at, name, address, unit_number, admin_id
		FROM apartments
		WHERE deleted_at IS NULL
	`

	// Basic filtering by ID or AdminID (extend this logic depending on your filter struct)
	args := []interface{}{}
	if f.ID != common.NilID {
		query += " AND id = $1"
		args = append(args, f.ID)
	} else {
		return nil, errors.New("no valid filter provided")
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	var apt types.Apartment
	err := row.Scan(
		&apt.ID,
		&apt.CreateAt,
		&apt.UpdateAt,
		&apt.DeleteAt,
		&apt.Name,
		&apt.Address,
		&apt.UnitNumber,
		&apt.AdminID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // apartment not found
		}
		return nil, err
	}

	return types.ApartmentStorageToDomain(&apt), nil
}

func (r *apartmentRepo) InviteMember(
	ctx context.Context,
	apartmentID common.ID,
	invite *domain.Invite,
) (
	*domain.Invite, error,
) {
	log := appctx.Logger(ctx)

	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO apartment_invites(apartment_id, invite_email, invite_status, invite_token, invite_expires_at)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id;`,
		apartmentID.String(), invite.Email,
		invite.Status, invite.Token, invite.ExpiresAt,
	).Scan(&id)
	if err != nil {
		log.Error("failed to execute query", zap.Error(err))
		return nil, err
	}

	err = invite.ID.UnmarshalText([]byte(id))
	if err != nil {
		log.Error("failed to unmarshal invite id", zap.Error(err))
		return nil, err
	}
	return invite, nil
}
