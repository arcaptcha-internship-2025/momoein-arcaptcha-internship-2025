package storage

import (
	"context"
	"database/sql"
	"errors"

	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/storage/types"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidFilter     = errors.New("no valid filter provided")
) 
type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) userPort.Repo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, ud *userDomain.User) (*userDomain.User, error) {
	log := appctx.Logger(ctx)

	u := types.UserDomainToStorage(ud)

	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users(email, password)
		VALUES($1, $2)
		ON CONFLICT (email) DO NOTHING
		RETURNING id;`,
		u.Email, u.Password,
	).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to execute query", zap.Error(ErrUserAlreadyExists))
			return nil, ErrUserAlreadyExists
		}
		log.Error("failed to execute query", zap.Error(err))
		return nil, err
	}

	u.ID = id
	return types.UserStorageToDomain(u), nil
}

func (r *userRepo) Get(
	ctx context.Context, filter *userDomain.UserFilter,
) (
	*userDomain.User, error,
) {
	log := appctx.Logger(ctx)

	var (
		query string
		args  []any
	)

	switch {
	case filter.ID != userDomain.NilID && filter.Email != "":
		query = `SELECT id, email, password, first_name, last_name FROM users WHERE id = $1 AND email = $2;`
		args = []any{filter.ID, filter.Email}
	case filter.ID != userDomain.NilID:
		query = `SELECT id, email, password, first_name, last_name FROM users WHERE id = $1;`
		args = []any{filter.ID}
	case filter.Email != "":
		query = `SELECT id, email, password, first_name, last_name FROM users WHERE email = $1;`
		args = []any{filter.Email}
	default:
		return nil, errors.New("no valid filter provided")
	}

	var u types.User
	err := r.db.QueryRowContext(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
	if err != nil {
		log.Error("failed to query user", zap.Error(err))
		return nil, err
	}

	user := types.UserStorageToDomain(&u)
	return user, nil
}

func (r *userRepo) Delete(ctx context.Context, f *userDomain.UserFilter) error {
	log := appctx.Logger(ctx)

	var (
		query string
		args  []any
	)

	if f.ID != userDomain.NilID && f.Email != "" {
		query = `DELETE FROM users WHERE id = $1 AND email = $2;`
		args = []any{f.ID, f.Email}
	} else if f.ID != userDomain.NilID {
		query = `DELETE FROM users WHERE id = $1;`
		args = []any{f.ID}
	} else if f.Email != "" {
		query = `DELETE FROM users WHERE email = $1;`
		args = []any{f.Email}
	} else {
		return ErrInvalidFilter
	}

	_, err := r.db.ExecContext(ctx, query, args...)

	if err != nil {
		log.Error("failed to execute query", zap.Error(err))
		return err
	}

	log.Info("deleting user",
		zap.String("id", f.ID.String()),
		zap.String("email", f.Email.String()),
	)
	return nil
}
