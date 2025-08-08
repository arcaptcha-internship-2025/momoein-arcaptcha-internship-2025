package paygw

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	paymentd "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	paymentp "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/port"
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

	callbackURL, err := url.Parse(tx.CallbackURL)
	if err != nil {
		return nil, err
	}
	if _, err = url.ParseQuery(callbackURL.RawQuery); err != nil {
		return nil, err
	}

	body, err := makeMapBody(&dto.PayRequest{
		Amount:      tx.Amount,
		CallbackURL: callbackURL.String(),
	})
	if err != nil {
		return nil, err
	}
	return &paymentd.RedirectGateway{
		Method: http.MethodPost,
		URL:    gatewayURL.String(),
		Body:   body,
	}, nil
}

func makeMapBody(body any) (map[string]any, error) {
	byteBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var mapBody map[string]any
	if err = json.Unmarshal(byteBody, &mapBody); err != nil {
		return nil, err
	}
	return mapBody, nil
}

func (g *mockGateway) VerifyTransaction(
	ctx context.Context,
	data map[string][]string,
) error {
	token, ok := data["token"]
	if !ok {
		return ErrMissingToken
	}
	verifyURL := *g.gatewayBaseURL
	verifyURL.Path = "/api/v1/payment/mock-gateway/verify"
	verifyURL.RawQuery = url.Values{"token": token}.Encode()

	resp, err := http.Get(verifyURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respBody dto.VerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}
	if respBody.Code != 0 {
		return ErrPaymentNotComplete
	}

	return nil
}
