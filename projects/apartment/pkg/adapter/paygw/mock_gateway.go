package paygw

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler"
	paymentd "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	paymentp "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/port"
)

const (
	PaymentIDsKey = "payment-ids"
)

var (
	ErrMissingToken       = errors.New("missing token")
	ErrPaymentNotComplete = errors.New("payment not complete")
)

type mockGateway struct {
	gatewayBaseURL *url.URL
}

func MustNewMockGateway(GatewayBaseURL string) paymentp.Gateway {
	gateway, err := NewMockGateway(GatewayBaseURL)
	if err != nil {
		panic(err)
	}
	return gateway
}

func NewMockGateway(GatewayBaseURL string) (paymentp.Gateway, error) {
	gbu, err := url.Parse(GatewayBaseURL)
	if err != nil {
		return nil, err
	}

	return &mockGateway{
		gatewayBaseURL: gbu,
	}, nil
}

func (g *mockGateway) CreateTransaction(
	ctx context.Context,
	tx paymentd.Transaction,
) (
	*paymentd.RedirectGateway, error,
) {
	gatewayURL := *g.gatewayBaseURL
	gatewayURL.Path = "/api/v1/payment/mock-gateway/pay"

	if _, err := url.Parse(tx.CallbackURL); err != nil {
		return nil, err
	}
	query := url.Values{}
	for i := range tx.PaymentIDs {
		query.Add(PaymentIDsKey, tx.PaymentIDs[i].String())
	}
	body := handler.PayRequest{
		Amount:    tx.Amount,
		ReturnURL: tx.CallbackURL,
		Query:     query.Encode(),
	}
	bBody, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	return &paymentd.RedirectGateway{
		Method: http.MethodPost,
		URL:    gatewayURL.String(),
		Body:   bBody,
	}, nil
}

func (g *mockGateway) VerifyTransaction(
	ctx context.Context,
	data map[string]string,
) error {
	token, ok := data["token"]
	if !ok {
		return ErrMissingToken
	}
	verifyURL := *g.gatewayBaseURL
	verifyURL.Path = "/api/v1/payment/mock-gateway/verify"
	verifyURL.RawQuery = url.Values{"token": {token}}.Encode()

	resp, err := http.Get(verifyURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respBody handler.VerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}
	if respBody.Code != 0 {
		return ErrPaymentNotComplete
	}

	return nil
}
