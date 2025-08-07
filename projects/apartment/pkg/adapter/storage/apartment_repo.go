package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment"
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	ap := types.ApartmentDomainToStorage(a)
	err = tx.QueryRowContext(ctx, `
		INSERT INTO apartments(name, address, unit_number, admin_id)
		VALUES($1, $2, $3, $4)
		RETURNING id;`,
		ap.Name, ap.Address, ap.UnitNumber, ap.AdminID,
	).Scan(&ap.ID)

	if err != nil {
		log.Error("failed to execute query", zap.Error(err))
		return nil, err
	}

	a = types.ApartmentStorageToDomain(ap)
	if err = r.AddUserToApartment(ctx, a.ID, a.AdminID); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *apartmentRepo) Get(
	ctx context.Context,
	f *domain.ApartmentFilter,
) (
	*domain.Apartment, error,
) {
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

func (r *apartmentRepo) AcceptInvite(ctx context.Context, token string) error {
	query := `
		SELECT invite_email, invite_status, invite_expires_at, apartment_id
		FROM apartment_invites
		WHERE invite_token = $1;
	`
	var (
		email  string
		status string
		exp    sql.NullTime
		aptId  string
	)
	err := r.db.QueryRowContext(ctx, query, token).Scan(&email, &status, &exp, &aptId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apartment.ErrInvalidToken
		}
		return err
	}
	if !exp.Valid || time.Since(exp.Time) > 0 {
		return apartment.ErrExpiredToken
	}
	if status == domain.InviteStatusAccepted.String() {
		return nil
	}

	var userId string
	err = r.db.QueryRowContext(ctx, `
		SELECT id FROM users WHERE email = $1 AND deleted_at IS NULL;
	`, email).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apartment.ErrUnregisteredUser
		}
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users_apartments(user_id, apartment_id)
		VALUES($1, $2);`, userId, aptId,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE apartment_invites
		SET invite_status = $1
		WHERE invite_token = $2;`, domain.InviteStatusAccepted.String(), token,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *apartmentRepo) AddUserToApartment(
	ctx context.Context, userId, aptId common.ID,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Check if user exists
	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, userId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user does not exist")
	}

	// Check if apartment exists
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM apartments WHERE id = $1)`, aptId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("apartment does not exist")
	}

	// Insert association
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users_apartments(user_id, apartment_id)
		VALUES ($1, $2);
	`, userId, aptId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
