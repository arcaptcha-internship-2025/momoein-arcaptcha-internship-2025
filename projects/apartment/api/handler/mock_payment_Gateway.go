package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"go.uber.org/zap"
)

const mockToken = "mock-token"

// MockGatewayPay
//
// @Summary      Mock payment gateway pay
// @Description  Simulates payment gateway pay endpoint for testing
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        body  body      dto.PayRequest  true  "Mock Payment Request"
// @Success      200   {object}  dto.PayResponse
// @Failure      400   {object}  dto.Error
// @Failure      502   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/payment/mock-gateway/pay [post]
func MockGatewayPay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		logPrefix := "mock-gateway-pay"

		var req dto.PayRequest
		if err := BodyParse(r, &req); err != nil {
			BadRequestError(w, r, "invalid request body")
			log.Error(fmt.Sprintf("%s: BodyParse", logPrefix), zap.Error(err))
			return
		}

		callbackURL, err := ParseURL(req.CallbackURL)
		if err != nil {
			BadRequestError(w, r, "invalid callback url")
			log.Error(fmt.Sprintf("%s: ParseURL", logPrefix), zap.Error(err))
			return
		}

		// Append mock token to query
		query := callbackURL.Query()
		query.Add("token", mockToken)
		callbackURL.RawQuery = query.Encode()

		// Perform callback request
		callBackResp, err := http.Post(callbackURL.String(), "application/json", nil)
		if err != nil {
			InternalServerError(w, r)
			log.Error(fmt.Sprintf("%s: CallBack", logPrefix), zap.Error(err))
			return
		}
		defer callBackResp.Body.Close()

		// Handle non-2xx callback responses
		if callBackResp.StatusCode >= 400 {
			log.Error(fmt.Sprintf("%s: callback returned error status: %d", logPrefix, callBackResp.StatusCode))
			callbackResult := fmt.Sprint(ResponseBodyToMap(callBackResp.Body))
			Error(w, r, http.StatusBadGateway, "callback endpoint returned error", callbackResult)
			return
		}

		// Prepare structured JSON response
		resp := dto.PayResponse{
			Status:  "completed",
			Message: "payment successfully completed",
			Token:   mockToken,
		}

		if err := WriteJson(w, http.StatusOK, &resp); err != nil {
			InternalServerError(w, r)
			log.Error(fmt.Sprintf("%s: WriteJson error", logPrefix), zap.Error(err))
		}
	})
}

func ResponseBodyToMap(body io.ReadCloser) map[string]any {
	content, _ := io.ReadAll(body)
	var contentMap map[string]any
	_ = json.Unmarshal(content, &contentMap)
	return contentMap
}

// MockGatewayVerify
//
// @Summary      Mock payment gateway verify
// @Description  Simulates payment gateway verify endpoint for testing
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        token  query    string  false "Mock Payment Token"
// @Param        body   body     dto.VerifyRequest false "Mock Verify Request"
// @Success      200   {object}  dto.VerifyResponse
// @Failure      400   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/payment/mock-gateway/verify [get]
func MockGatewayVerify() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		logPrefix := "mock-gateway-verify"

		var req dto.VerifyRequest
		if r.URL.Query().Has("token") {
			req.Token = r.URL.Query().Get("token")
		} else {
			if err := BodyParse(r, &req); err != nil {
				BadRequestError(w, r, "invalid request body")
				log.Error(fmt.Sprintf("%s: BodyParse", logPrefix), zap.Error(err))
				return
			}
		}

		resp := dto.VerifyResponse{
			Code:    0,
			Message: "payment completed",
		}

		if req.Token != mockToken {
			resp.Code = 1
			resp.Message = "payment not complete"
		}

		if err := WriteJson(w, http.StatusOK, &resp); err != nil {
			InternalServerError(w, r)
			log.Error(fmt.Sprintf("%s: WriteJson", logPrefix), zap.Error(err))
		}
	})
}

func ParseURL(URL string) (u *url.URL, err error) {
	u, err = url.Parse(URL)
	if err != nil {
		return nil, err
	}
	_, err = url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	return
}
