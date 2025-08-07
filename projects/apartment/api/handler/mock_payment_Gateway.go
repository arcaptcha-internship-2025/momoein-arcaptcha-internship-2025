package handler

import (
	"fmt"
	"net/http"
	"net/url"

	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"go.uber.org/zap"
)

const mockToken = "mock-token"

type PayRequest struct {
	Amount    int64  `json:"amount"`
	ReturnURL string `json:"returnUrl"`
	Query     string `json:"query"` // Query is expected to be a list of key=value settings separated by ampersands.
}

type PayResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type VerifyResponse struct {
	Code    int    `json:"code"`    //  ==0 success, !=0 failed
	Message string `json:"message"` // descriptive message
}

func MockGatewayPay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		logPrefix := "mock-gateway-pay"

		var req PayRequest
		if err := BodyParse(r, &req); err != nil {
			BadRequestError(w, r, "invalid request body")
			log.Error(fmt.Sprintf("%s: BodyParse", logPrefix), zap.Error(err))
			return
		}

		returnURL, err := ParseURL(req.ReturnURL, req.Query)
		if err != nil {
			BadRequestError(w, r, "invalid callback url")
			log.Error(fmt.Sprintf("%s: ParseURL", logPrefix), zap.Error(err))
			return
		}

		// Append mock token to query
		query := returnURL.Query()
		query.Add("token", mockToken)
		returnURL.RawQuery = query.Encode()

		// Perform callback request
		callBackResp, err := http.Post(returnURL.String(), "application/json", nil)
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
		resp := PayResponse{
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

		var req VerifyRequest
		if r.URL.Query().Has("token") {
			req.Token = r.URL.Query().Get("token")
		} else {
			if err := BodyParse(r, &req); err != nil {
				BadRequestError(w, r, "invalid request body")
				log.Error(fmt.Sprintf("%s: BodyParse", logPrefix), zap.Error(err))
				return
			}
		}

		resp := VerifyResponse{
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

func ParseURL(
	returnURL, query string,
) (
	ru *url.URL, err error,
) {
	ru, err = url.Parse(returnURL)
	if err != nil {
		return nil, err
	}
	_, err = url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	ru.RawQuery = query
	return
}
