package payment

import (
	"context"
	"errors"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
)

var (
	ErrOnPayBill       = errors.New("error on pay bill")
	ErrNoBalanceDue    = errors.New("no balance due")
	ErrUnknownGateway  = errors.New("unknown gateway")
	ErrOnPayTotalDebt  = errors.New("error on pay total debt")
	ErrOnCallback      = errors.New("error on handle callbackI")
	ErrInvalidCallback = errors.New("invalid callback")
)

type service struct {
	repo     port.Repo
	gateways map[domain.GatewayType]port.Gateway
}

func NewService(repo port.Repo, gws map[domain.GatewayType]port.Gateway) port.Service {
	return &service{
		repo:     repo,
		gateways: gws,
	}
}

func (s *service) Gateway(gt domain.GatewayType) (port.Gateway, error) {
	gateway, ok := s.gateways[gt]
	if !gt.IsValid() || !ok {
		return nil, ErrUnknownGateway
	}
	return gateway, nil
}

func (s *service) PayBill(
	ctx context.Context,
	gt domain.GatewayType,
	userID, billID common.ID,
	callBackURL string,
) (
	redirectURL string, err error,
) {
	gateway, err := s.Gateway(gt)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayBill, err)
	}

	balanceDue, err := s.repo.UserBillBalanceDue(ctx, userID, billID)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayBill, err)
	}
	if balanceDue <= 0 {
		return "", fp.WrapErrors(ErrOnPayBill, ErrNoBalanceDue)
	}

	p := &domain.Payment{BillID: billID, PayerID: userID, Amount: balanceDue}

	_, err = s.repo.CreatePayment(ctx, p)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayBill, err)
	}

	tx := domain.Transaction{
		Amount:  balanceDue,
		PayerID: userID,
		Bills: []domain.BillWithAmount{{
			BillID: billID,
			Amount: balanceDue,
		}},
		CallbackURL: callBackURL,
		Metadata:    map[string]string{"payment_type": "total-debt"},
	}

	redirectURL, err = gateway.CreateTransaction(ctx, tx)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayBill, err)
	}
	return
}

func (s *service) PayTotalDebt(
	ctx context.Context,
	gt domain.GatewayType,
	userID common.ID,
	callBackURL string,
) (
	redirectURL string, err error,
) {
	gateway, err := s.Gateway(gt)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayTotalDebt, err)
	}

	balanceDues, err := s.repo.UserBillsBalanceDue(ctx, userID)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayTotalDebt, err)
	}
	if len(balanceDues) == 0 {
		return "", fp.WrapErrors(ErrOnPayTotalDebt, ErrNoBalanceDue)
	}

	var totalAmount int64
	var payments []*domain.Payment

	for _, bDue := range balanceDues {
		if bDue.Amount <= 0 {
			continue
		}
		p := &domain.Payment{
			PayerID: userID,
			BillID:  bDue.BillID,
			Amount:  bDue.Amount,
			Status:  domain.PaymentPending,
			Gateway: gt.String(),
		}
		payments = append(payments, p)
		totalAmount += bDue.Amount
	}

	_, err = s.repo.BatchCreatePayment(ctx, payments)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayTotalDebt, err)
	}

	tx := domain.Transaction{
		Amount:      totalAmount,
		PayerID:     userID,
		Bills:       balanceDues,
		CallbackURL: "",
		Metadata:    map[string]string{"payment_type": "total-debt"},
	}

	redirectURL, err = gateway.CreateTransaction(ctx, tx)
	if err != nil {
		return "", fp.WrapErrors(ErrOnPayTotalDebt, err)
	}
	return
}

func (s *service) HandleCallback(
	ctx context.Context,
	gt domain.GatewayType,
	data map[string]string,
) error {
	gateway, err := s.Gateway(gt)
	if err != nil {
		return fp.WrapErrors(ErrOnCallback, err)
	}

	if err := gateway.VerifyTransaction(ctx, data); err != nil {
		return fp.WrapErrors(ErrOnCallback, ErrInvalidCallback, err)
	}

	// TODO
	// change status if payment successful
	return nil
}
