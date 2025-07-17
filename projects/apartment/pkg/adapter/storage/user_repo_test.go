package storage

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/storage/types"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_Create(t *testing.T) {
	// Step 1: Create mock DB and expectations
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Step 2: Prepare expected input values
	email := "momoein@gmail.com"
	password := "my-password"
	fakeUUID := uuid.New()

	// Step 3: Create a user domain object
	user := &userDomain.User{
		Email: common.Email(email),
	}
	user.SetPassword([]byte(password)) // sets internal password
	require.NoError(t, err)

	// Step 4: Convert to storage format to get expected password value
	storageUser := types.UserDomainToStorage(user)

	// Step 5: Mock the expected query and result
	mock.ExpectPrepare(`INSERT INTO user\(email, password\) VALUES\(\$1, \$2\) RETURNING id;`).
		ExpectQuery().
		WithArgs(storageUser.Email, storageUser.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(fakeUUID.String()))

	// Step 6: Create repo instance
	repo := NewUserRepo(db)

	// Step 7: Create context with logger
	ctx := appctx.New(context.Background(),
		appctx.WithLogger(logger.NewConsoleZapLogger(logger.ModeDevelopment)))

	// Step 8: Call Create method
	createdUser, err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.Equal(t, fakeUUID, createdUser.ID)
	require.Equal(t, user.Email, createdUser.Email)

	// Step 9: Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
