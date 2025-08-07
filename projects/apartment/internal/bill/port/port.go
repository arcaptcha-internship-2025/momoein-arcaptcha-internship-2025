package port

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
)

type Service interface {
	AddBill(context.Context, *domain.Bill) (*domain.Bill, error)
	GetBill(context.Context, *domain.BillFilter) (*domain.Bill, error)
	GetBillImage(ctx context.Context, imageID common.ID) (string, error)
	GetUserBillShares(ctx context.Context, userID common.ID) ([]domain.UserBillShare, error)
	GetUserTotalDebt(ctx context.Context, userID common.ID) (int, error)
}

type Repo interface {
	Create(context.Context, *domain.Bill) (*domain.Bill, error)
	Read(context.Context, *domain.BillFilter) (*domain.Bill, error)
	GetUserBillShares(ctx context.Context, userID common.ID) ([]domain.UserBillShare, error)
	GetUserTotalDebt(ctx context.Context, userID common.ID) (int, error)
}

type ObjectStorage interface {
	FPut(ctx context.Context, key, filename string) error
	FGet(ctx context.Context, key, filename string) error
	Del(ctx context.Context, key string) error
}
