package user

import (
	"context"
	"errors"
	"testing"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	log = logger.NewConsoleZapLogger(logger.ModeDevelopment)
	ctx = appctx.New(context.Background(), appctx.WithLogger(log))
)

type MockRepo struct {
	mock.Mock
	port.Repo
}

func (m *MockRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepo) Get(ctx context.Context, f *domain.UserFilter) (*domain.User, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepo) Delete(ctx context.Context, f *domain.UserFilter) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	u := domain.NewUser(
		common.NewRandomID(),
		"test@gmail.com",
		"password",
		"", "",
	)

	repo.On("Create", ctx, u).Return(u, nil)

	usr, err := svc.Create(ctx, u)

	assert.NoError(t, err)
	assert.NotNil(t, usr)

	repo.AssertExpectations(t)
}

func TestCreate_ValidationError(t *testing.T) {
	testData := []struct {
		user *domain.User
		err  error
	}{
		{
			user: &domain.User{},
			err:  domain.ErrInvalidEmail,
		},
		{
			user: &domain.User{Email: "test@gmail.com"},
			err:  domain.ErrUserShortPassword,
		},
		{
			user: domain.NewUser(common.NilID, "test@gmail.com", string(make([]byte, 73)), "", ""),
			err:  domain.ErrUserLongPassword,
		},
	}

	repo := new(MockRepo)
	svc := NewService(repo)

	for _, test := range testData {
		u := test.user
		err := u.Validate()
		assert.Error(t, err)

		usr, err := svc.Create(ctx, u)

		assert.Nil(t, usr)
		assert.Error(t, err)
		assert.ErrorIs(t, err, test.err)
		assert.Contains(t, err.Error(), test.err.Error())
	}
}

func TestCreate_RepoError(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	u := domain.NewUser(common.NewRandomID(), "test@gmail.com", "password", "", "")

	repo.On("Create", ctx, u).Return((*domain.User)(nil), errors.New("db error"))

	usr, err := svc.Create(ctx, u)

	assert.Error(t, err)
	assert.Nil(t, usr)
	assert.Contains(t, err.Error(), "user Creation failed")
	repo.AssertExpectations(t)
}

func TestGet_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	user := domain.NewUser(common.NewRandomID(), "test@gmail.com", "password", "", "")
	filter := &domain.UserFilter{Email: "test@gmail.com"}

	repo.On("Get", ctx, filter).Return(user, nil)

	u, err := svc.Get(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, user.Email, u.Email)
	repo.AssertExpectations(t)
}

func TestGet_InvalidFilter(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	filter := &domain.UserFilter{}

	u, err := svc.Get(ctx, filter)

	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidOrNilFilter, err)
}

func TestGet_RepoError(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	filter := &domain.UserFilter{Email: "test@gmail.com"}

	repo.On("Get", ctx, filter).Return((*domain.User)(nil), errors.New("db error"))

	u, err := svc.Get(ctx, filter)

	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Contains(t, err.Error(), "user retrieve failed")
	repo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	filter := &domain.UserFilter{Email: "test@gmail.com"}

	repo.On("Delete", ctx, filter).Return(nil)

	err := svc.Delete(ctx, filter)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDelete_InvalidFilter(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	filter := &domain.UserFilter{}

	err := svc.Delete(ctx, filter)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidOrNilFilter, err)
}

func TestDelete_RepoError(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	ctx := context.Background()
	filter := &domain.UserFilter{Email: "test@gmail.com"}

	repo.On("Delete", ctx, filter).Return(errors.New("db error"))

	err := svc.Delete(ctx, filter)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error on deleting failed")
	repo.AssertExpectations(t)
}
