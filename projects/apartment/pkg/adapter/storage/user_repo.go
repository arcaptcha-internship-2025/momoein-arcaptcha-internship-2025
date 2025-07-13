package storage

import (
	"context"
	"database/sql"

	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/storage/types"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"go.uber.org/zap"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) userPort.Repo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, ud *userDomain.User) (*userDomain.User, error) {
	// Convert domain model to storage model
	u := types.UserDomainToStorage(ud)

	// Prepare the SQL statement with RETURNING id
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO user(email, password) VALUES($1, $2) RETURNING id;`)
	if err != nil {
		appctx.Logger(ctx).Error("failed to prepare insert user statement", zap.Error(err))
		return nil, err
	}
	defer stmt.Close()

	// Get the auto-generated ID
	var id string
	err = stmt.QueryRowContext(ctx, u.Email, u.Password).Scan(&id)
	if err != nil {
		appctx.Logger(ctx).Error("failed to execute insert statement", zap.Error(err))
		return nil, err
	}

	u.ID = id
	return types.UserStorageToDomain(u), nil
}



func (r *userRepo) Get(context.Context, *userDomain.UserFilter) (*userDomain.User, error) {
	panic("unimplemented")
}
func (r *userRepo) Delete(context.Context, *userDomain.UserFilter) error {
	panic("unimplemented")
}
