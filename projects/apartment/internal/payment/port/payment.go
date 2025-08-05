package port

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
)

type Service interface {
	PayBill(ctx context.Context, gateway domain.GatewayType, billID, userID common.ID) (redirectURL string, err error)
	PayTotalDebt(ctx context.Context, gateway domain.GatewayType, userID common.ID) (redirectURL string, err error)
	HandleCallback(ctx context.Context, gateway domain.GatewayType, data map[string]string) error
}

type Repo interface {
	CreatePayment(ctx context.Context, p *domain.Payment) (*domain.Payment, error)
	UpdateStatus(ctx context.Context, paymentID common.ID, s domain.PaymentStatus) error
	UserBillBalanceDue(ctx context.Context, userId, billId common.ID) (int64, error)
}

type Gateway interface {
	CreateTransaction(ctx context.Context, payment *domain.Payment) (redirectURL string, err error)
	VerifyTransaction(ctx context.Context, data map[string]string) (bool, error)
}
