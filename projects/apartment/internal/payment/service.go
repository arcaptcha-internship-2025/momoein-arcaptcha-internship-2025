package payment

import (
	"context"
	"errors"
	"net/url"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
)

const (
	PaymentIDsKey = "payment-ids"
)

var (
	ErrOnPayBill       = errors.New("error on pay bill")
	ErrNoBalanceDue    = errors.New("no balance due")
	ErrUnknownGateway  = errors.New("unknown gateway")
	ErrOnPayTotalDebt  = errors.New("error on pay total debt")
	ErrOnCallback      = errors.New("error on handle callbackI")
	ErrInvalidCallback = errors.New("invalid callback")
	ErrInvalidStatus   = errors.New("invalid status")
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
	redirect *domain.RedirectGateway, err error,
) {
	gateway, err := s.Gateway(gt)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayBill, err)
	}
	balanceDue, err := s.repo.UserBillBalanceDue(ctx, userID, billID)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayBill, err)
	}
	if balanceDue <= 0 {
		return nil, fp.WrapErrors(ErrOnPayBill, ErrNoBalanceDue)
	}
	p := &domain.Payment{BillID: billID, PayerID: userID, Amount: balanceDue}
	p, err = s.repo.CreatePayment(ctx, p)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayBill, err)
	}
	callBackURL, err = CallbackURLWithPaymentIDs(callBackURL, p.ID)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayBill, err)
	}
	tx := domain.Transaction{
		Amount:  balanceDue,
		PayerID: userID,
		Bills: []domain.BillWithAmount{{
			BillID: billID,
			Amount: balanceDue,
		}},
		CallbackURL: callBackURL,
	}
	redirect, err = gateway.CreateTransaction(ctx, tx)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayBill, err)
	}
	return
}

func (s *service) PayTotalDebt(
	ctx context.Context,
	gt domain.GatewayType,
	userID common.ID,
	callBackURL string,
) (
	redirect *domain.RedirectGateway, err error,
) {
	gateway, err := s.Gateway(gt)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayTotalDebt, err)
	}
	balanceDues, err := s.repo.UserBillsBalanceDue(ctx, userID)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayTotalDebt, err)
	}
	if len(balanceDues) == 0 {
		return nil, fp.WrapErrors(ErrOnPayTotalDebt, ErrNoBalanceDue)
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
	payments, err = s.repo.BatchCreatePayment(ctx, payments)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayTotalDebt, err)
	}
	paymentIDs := []common.ID{}
	for i := range payments {
		paymentIDs = append(paymentIDs, payments[i].ID)
	}
	callBackURL, err = CallbackURLWithPaymentIDs(callBackURL, paymentIDs...)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayBill, err)
	}
	tx := domain.Transaction{
		PaymentIDs:  paymentIDs,
		Amount:      totalAmount,
		PayerID:     userID,
		Bills:       balanceDues,
		CallbackURL: callBackURL,
	}
	redirect, err = gateway.CreateTransaction(ctx, tx)
	if err != nil {
		return nil, fp.WrapErrors(ErrOnPayTotalDebt, err)
	}
	return
}

func CallbackURLWithPaymentIDs(callbackURL string, IDs ...common.ID) (string, error) {
	cURL, err := url.Parse(callbackURL)
	if err != nil {
		return "", err
	}
	query, err := url.ParseQuery(cURL.RawQuery)
	if err != nil {
		return "", err
	}
	for i := range IDs {
		query.Add(PaymentIDsKey, IDs[i].String())
	}
	cURL.RawQuery = query.Encode()
	return cURL.String(), nil
}

func (s *service) HandleCallback(
	ctx context.Context,
	gt domain.GatewayType,
	data map[string][]string,
) error {
	gateway, err := s.Gateway(gt)
	if err != nil {
		return fp.WrapErrors(ErrOnCallback, err)
	}
	err = gateway.VerifyTransaction(ctx, data)
	if err != nil {
		return fp.WrapErrors(ErrOnCallback, ErrInvalidCallback, err)
	}
	paymentIDs := []common.ID{}
	for _, id := range data[PaymentIDsKey] {
		paymentIDs = append(paymentIDs, common.IDFromText(id))
	}
	err = s.repo.UpdateStatus(ctx, paymentIDs, domain.PaymentPaid)
	if err != nil {
		return fp.WrapErrors(ErrOnCallback, err)
	}
	return nil
}

func (s *service) SupportedGateways() []string {
	result := make([]string, 0)
	for key := range s.gateways {
		result = append(result, key.String())
	}
	return result
}
