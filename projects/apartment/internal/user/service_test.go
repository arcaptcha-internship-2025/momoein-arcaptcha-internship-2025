package user

import (
	"context"
	"testing"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	return args.Error(1)
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
