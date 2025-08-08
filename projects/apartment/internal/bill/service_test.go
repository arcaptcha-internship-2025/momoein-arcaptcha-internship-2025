package bill

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	log = logger.NewConsoleZapLogger(logger.ModeDevelopment)
	ctx = appctx.New(context.Background(), appctx.WithLogger(log))
)

// ----------- Mocks -------------

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, b *domain.Bill) (*domain.Bill, error) {
	args := m.Called(ctx, b)
	return args.Get(0).(*domain.Bill), args.Error(1)
}

func (m *MockRepo) Read(ctx context.Context, f *domain.BillFilter) (*domain.Bill, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(*domain.Bill), args.Error(1)
}

func (m *MockRepo) GetUserBillShares(ctx context.Context, userID common.ID) ([]domain.UserBillShare, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.UserBillShare), args.Error(1)
}

func (m *MockRepo) GetUserTotalDebt(ctx context.Context, userID common.ID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Set(key string, val any) error {
	args := m.Called(key, val)
	return args.Error(0)
}

func (m *MockStorage) Get(key string) any {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockStorage) FPut(ctx context.Context, key, filename string) error {
	args := m.Called(ctx, key, filename)
	return args.Error(0)
}

func (m *MockStorage) FGet(ctx context.Context, key, filename string) error {
	args := m.Called(ctx, key, filename)
	return args.Error(0)
}

func (m *MockStorage) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func createValidBill() *domain.Bill {
	return &domain.Bill{
		ID:          common.NewRandomID(),
		Type:        domain.BillWater,
		BillNumber:  2345,
		Amount:      100,
		HasImage:    true,
		Image:       &domain.Image{Name: "mockImageData"},
		DueDate:     time.Now().Add(3 * 24 * time.Hour),
		ApartmentID: common.NewRandomID(),
	}
}

func createValidFilter() *domain.BillFilter {
	return &domain.BillFilter{
		ID: common.NewRandomID(),
	}
}

// ----------- Tests -------------

func TestAddBill_SuccessWithImage(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	bill := createValidBill()

	storage.On("FPut", ctx, mock.Anything, mock.Anything).Return(nil)
	repo.On("Create", ctx, bill).Return(bill, nil)

	result, err := svc.AddBill(ctx, bill)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestAddBill_ValidationError(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	tests := []struct {
		bill *domain.Bill
		err  error
	}{
		{
			bill: createValidBill().SetType(""),
			err:  domain.ErrBillMissingType,
		},
		{
			bill: createValidBill().SetBillNumber(0),
			err:  domain.ErrBillInvalidBillNumber,
		},
		{
			bill: createValidBill().SetAmount(-1),
			err:  domain.ErrBillNegativeAmount,
		},
		{
			bill: createValidBill().SetDueDate(time.Time{}),
			err:  domain.ErrBillMissingDueDate,
		},
	}

	for _, test := range tests {
		bill := test.bill
		err := bill.Validate()
		assert.Error(t, err)

		result, err := svc.AddBill(ctx, bill)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid bill")
	}
}

func TestAddBill_ImageStorageError(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	bill := createValidBill()

	storage.On("FPut", ctx, mock.Anything, mock.Anything).Return(errors.New("storage failed"))

	result, err := svc.AddBill(ctx, bill)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error on add bill")
	repo.AssertNotCalled(t, "Create")
}

func TestGetBill_SuccessWithImage(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	filter := createValidFilter()
	expectedBill := createValidBill()

	repo.On("Read", ctx, filter).Return(expectedBill, nil)
	storage.On("FGet", ctx, mock.Anything, mock.Anything).Return(nil)

	result, err := svc.GetBill(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBill.Image, result.Image)
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestGetBill_SuccessBadImageFormat(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	filter := createValidFilter()
	bill := createValidBill()
	bill.Image = nil

	repo.On("Read", ctx, filter).Return(bill, nil)
	storage.On("FGet", ctx, mock.Anything, mock.Anything).
		Return(errors.New("failed to get image"))

	result, err := svc.GetBill(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.Image) // image shouldn't be set
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestGetBill_RepoError(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	filter := createValidFilter()
	bill := createValidBill()
	bill.Image = nil

	repo.On("Read", ctx, filter).Return(bill, errors.New("repo error"))
	storage.On("Get", bill.ImageID.String()).Return(domain.Image{})

	result, err := svc.GetBill(ctx, filter)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), ErrOnGetBill.Error())
	repo.AssertExpectations(t)
}

func TestGetBillImage_Success(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	imageID := common.NewRandomID()
	expectedPath := filepath.Join(os.TempDir(), imageID.String())

	storage.
		On("FGet", mock.Anything, imageID.String(), expectedPath).
		Return(nil)

	path, err := svc.GetBillImage(context.Background(), imageID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, path)
	storage.AssertExpectations(t)
}

func TestGetBillImage_ObjectStorageError(t *testing.T) {
	repo := new(MockRepo)
	storage := new(MockStorage)
	svc := NewService(repo, storage)

	imageID := common.NewRandomID()
	expectedPath := filepath.Join(os.TempDir(), imageID.String())
	storageErr := errors.New("storage unavailable")

	storage.
		On("FGet", mock.Anything, imageID.String(), expectedPath).
		Return(storageErr)

	path, err := svc.GetBillImage(context.Background(), imageID)

	assert.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "storage unavailable")
	storage.AssertExpectations(t)
}

func TestGetBillImage_NotFound(t *testing.T) {
	// TODO
}
