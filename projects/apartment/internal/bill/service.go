package bill

import (
	"context"
	"errors"
	"image"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
)

var (
	ErrOnAddBill = errors.New("error on add bill")
	ErrOnGetBill = errors.New("error on get bill")
	ErrNotFound  = errors.New("source not found")
)

type service struct {
	repo port.Repo
	strg port.Storage
}

func NewService(r port.Repo, s port.Storage) port.Service {
	return &service{repo: r, strg: s}
}

func (s *service) AddBill(ctx context.Context, bill *domain.Bill) (*domain.Bill, error) {
	if bill.HasImage {
		bill.ImageID = common.NewRandomID()
		err := s.strg.Set(bill.ImageID.String(), bill.Image)
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
		img, ok := s.strg.Get(bill.ImageID.String()).(image.Image)
		if !ok {
			log.Warn("bad image format")
		} else {
			bill.Image = img
		}
	}

	return bill, nil
}
