package port

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
)

type Service interface {
	AddBill(context.Context, *domain.Bill) (*domain.Bill, error)
	GetBill(context.Context, *domain.BillFilter) (*domain.Bill, error)
}

type Repo interface {
	Create(context.Context, *domain.Bill) (*domain.Bill, error)
	Read(context.Context, *domain.BillFilter) (*domain.Bill, error)
}

type Storage interface {
	Set(key string, val any) error
	Get(key string) any
	Del(key string) error
}
