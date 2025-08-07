package apartment

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
)

type MockRepo struct {
	mock.Mock
	port.Repo
}

func (m *MockRepo) Create(
	ctx context.Context, a *domain.Apartment,
) (
	*domain.Apartment, error,
) {
	args := m.Called(ctx, a)
	return args.Get(0).(*domain.Apartment), args.Error(1)
}

func (m *MockRepo) Get(
	ctx context.Context, filter *domain.ApartmentFilter,
) (
	*domain.Apartment, error,
) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*domain.Apartment), args.Error(1)
}

func (m *MockRepo) InviteMember(
	ctx context.Context, aptID common.ID, invite *domain.Invite,
) (
	*domain.Invite, error,
) {
	args := m.Called(ctx, aptID, invite)
	return args.Get(0).(*domain.Invite), args.Error(1)
}

func (m *MockRepo) AcceptInvite(
	ctx context.Context, token string,
) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

type MockEmail struct {
	mock.Mock
	port.EmailSender
}

func (m *MockEmail) Send(to []string, msg *common.EmailMessage) error {
	args := m.Called(to, msg)
	return args.Error(0)
}

var (
	log = logger.NewConsoleZapLogger(logger.ModeDevelopment)
	ctx = appctx.New(context.Background(), appctx.WithLogger(log))
)

func getAcceptURL() url.URL {
	return url.URL{
		Scheme: "http",
		Host:   "127.0.0.1:8080",
		Path:   "api/v1/apartment/invite/accept",
	}
}

func TestCreateApartment_Success(t *testing.T) {
	repo := new(MockRepo)
	email := new(MockEmail)
	svc := NewService(repo, email)

	a := &domain.Apartment{
		ID:      common.NewRandomID(),
		Name:    "Test Apt",
		AdminID: common.NewRandomID(),
	}

	repo.On("Create", ctx, a).Return(a, nil)

	apt, err := svc.Create(ctx, a)

	assert.NoError(t, err)
	assert.Equal(t, a, apt)
	repo.AssertExpectations(t)
}

func TestInviteMember_Success(t *testing.T) {
	repo := new(MockRepo)
	email := new(MockEmail)
	svc := NewService(repo, email)

	adminID := common.NewRandomID()
	apartmentID := common.NewRandomID()
	userEmail := common.Email("test@example.com")
	acceptURL := getAcceptURL()

	apartmentObj := &domain.Apartment{
		ID:      apartmentID,
		AdminID: adminID,
	}

	invite := &domain.Invite{
		Email:     userEmail,
		Status:    domain.InviteStatusPending,
		Token:     common.NewRandomID().String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	repo.On("Get", ctx, &domain.ApartmentFilter{ID: apartmentID}).Return(apartmentObj, nil)
	repo.On("InviteMember", ctx, apartmentID, mock.AnythingOfType("*domain.Invite")).Return(invite, nil)
	email.On("Send", []string{userEmail.String()}, mock.Anything).Return(nil)

	result, err := svc.InviteMember(ctx, adminID, apartmentID, userEmail, acceptURL.String())

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, invite.Email, result.Email)
	}
	repo.AssertExpectations(t)
	email.AssertExpectations(t)
}
func TestInviteMember_NotAdmin(t *testing.T) {
	repo := new(MockRepo)
	email := new(MockEmail)
	svc := NewService(repo, email)

	adminID := common.NewRandomID()
	apartmentID := common.NewRandomID()
	userEmail := common.Email("user@example.com")
	acceptURL := getAcceptURL()

	apt := &domain.Apartment{
		ID:      apartmentID,
		AdminID: common.NewRandomID(), // different admin
	}

	repo.On("Get", ctx, &domain.ApartmentFilter{ID: apartmentID}).Return(apt, nil)

	invite, err := svc.InviteMember(ctx, adminID, apartmentID, userEmail, acceptURL.String())

	assert.Error(t, err)
	assert.Nil(t, invite)
	repo.AssertExpectations(t)
}

func TestAcceptInvite_Success(t *testing.T) {
	repo := new(MockRepo)
	email := new(MockEmail)
	svc := NewService(repo, email)

	token := common.NewRandomID().String()

	repo.On("AcceptInvite", ctx, token).Return(nil)

	err := svc.AcceptInvite(ctx, token)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
func TestAcceptInvite_InvalidToken(t *testing.T) {
	repo := new(MockRepo)
	email := new(MockEmail)
	svc := NewService(repo, email)

	token := "k3jd4kj9e-kdh4iu-ejf4ioj4k"

	err := svc.AcceptInvite(ctx, token)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestAcceptInvite_ExpiredToken(t *testing.T) {
	repo := new(MockRepo)
	email := new(MockEmail)
	svc := NewService(repo, email)

	token := common.NewRandomID().String()

	repo.On("AcceptInvite", ctx, token).Return(ErrExpiredToken)

	err := svc.AcceptInvite(ctx, token)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrExpiredToken)
	repo.AssertExpectations(t)

}

func TestAcceptInvite_Unregistered(t *testing.T) {
	repo := new(MockRepo)
	email := new(MockEmail)
	svc := NewService(repo, email)

	token := common.NewRandomID().String()

	repo.On("AcceptInvite", ctx, token).Return(ErrUnregisteredUser)

	err := svc.AcceptInvite(ctx, token)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnregisteredUser)
	repo.AssertExpectations(t)
}
