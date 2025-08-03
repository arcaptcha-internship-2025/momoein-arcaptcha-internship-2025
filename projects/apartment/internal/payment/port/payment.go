package port

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, billID string, gateway string) (*domain.Payment, error)
	HandleCallback(ctx context.Context, gateway string, data map[string]string) error
}
