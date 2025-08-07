package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"go.uber.org/zap"
)

const mockToken = "mock-token"

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
			Error(w, r, http.StatusBadGateway, "callback endpoint returned error")
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
