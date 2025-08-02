package bill

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
	"go.uber.org/zap"
)

var (
	ErrOnAddBill      = errors.New("error on add bill")
	ErrOnGetBill      = errors.New("error on get bill")
	ErrNotFound       = errors.New("source not found")
	ErrBillOnValidate = errors.New("invalid bill")
	ErrOnGetBillImage = errors.New("error on get bill image")
)

type service struct {
	repo port.Repo
	strg port.ObjectStorage
}

func NewService(r port.Repo, s port.ObjectStorage) port.Service {
	return &service{repo: r, strg: s}
}

func (s *service) AddBill(ctx context.Context, bill *domain.Bill) (*domain.Bill, error) {
	log := appctx.Logger(ctx)

	if err := bill.Validate(); err != nil {
		log.Error("service AddBill", zap.Error(err))
		return nil, fp.WrapErrors(ErrOnAddBill, ErrBillOnValidate, err)
	}

	if bill.HasImage && bill.Image != nil {
		bill.ImageID = common.NewRandomID()
		err := s.strg.FPut(ctx, bill.ImageID.String(), bill.Image.Path)
		if err != nil {
			return nil, fp.WrapErrors(ErrOnAddBill, err)
		}
	}
	return s.repo.Create(ctx, bill)
}

func (s *service) GetBill(ctx context.Context, f *domain.BillFilter) (*domain.Bill, error) {
	log := appctx.Logger(ctx)

	bill, err := s.repo.Read(ctx, f)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnGetBill, err)
	}

	if bill.HasImage {
		path, err := s.GetBillImage(ctx, bill.ImageID)
		if err != nil {
			log.Warn("failed to get bill image", zap.Error(err))
		} else {
			bill.Image = &domain.Image{Path: path}
		}
	}

	return bill, nil
}

func (s *service) GetBillImage(ctx context.Context, imageID common.ID) (string, error) {
	path := filepath.Join(os.TempDir(), imageID.String())
	err := s.strg.FGet(ctx, imageID.String(), path)
	if err != nil {
		return "", fp.WrapErrors(ErrOnGetBillImage, err)
	}
	return path, nil
}

func (s *service) GetUserBillShares(
	ctx context.Context, userID common.ID,
) (
	[]domain.UserBillShare, error,
) {
	ubs, err := s.repo.GetUserBillShares(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(ubs) == 0 {
		return nil, ErrNotFound
	}
	return ubs, nil
}

func (s *service) GetUserTotalDebt(ctx context.Context, userID common.ID) (int, error) {
	debt, err := s.repo.GetUserTotalDebt(ctx, userID)
	if err != nil {
		return 0, err
	}
	return debt, nil
}
