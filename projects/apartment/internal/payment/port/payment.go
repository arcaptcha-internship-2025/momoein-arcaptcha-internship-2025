package port

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
)

type Service interface {
	PayBill(ctx context.Context, gateway domain.GatewayType, billID, userID common.ID, callBackURL string) (*domain.RedirectGateway, error)
	PayTotalDebt(ctx context.Context, gateway domain.GatewayType, userID common.ID, callBackURL string) (*domain.RedirectGateway, error)
	HandleCallback(ctx context.Context, gateway domain.GatewayType, data map[string][]string) error
}

type Repo interface {
	CreatePayment(ctx context.Context, p *domain.Payment) (*domain.Payment, error)
	BatchCreatePayment(ctx context.Context, ps []*domain.Payment) ([]*domain.Payment, error)
	UpdateStatus(ctx context.Context, paymentID []common.ID, s domain.PaymentStatus) error
	UserBillBalanceDue(ctx context.Context, userId, billId common.ID) (int64, error)
	UserBillsBalanceDue(ctx context.Context, userId common.ID) ([]domain.BillWithAmount, error)
}

type Gateway interface {
	CreateTransaction(ctx context.Context, tx domain.Transaction) (*domain.RedirectGateway, error)
	VerifyTransaction(ctx context.Context, data map[string][]string) error
}
