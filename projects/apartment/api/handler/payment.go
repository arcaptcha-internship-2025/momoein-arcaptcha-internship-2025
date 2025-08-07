package handler

import (
	"errors"
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment"
	paymentd "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	paymentp "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/port"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	appjwt "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/jwt"
	"go.uber.org/zap"
)

// PayUserBill
//
// @Summary      Pay a bill
// @Description  Initiates payment for a specific bill using a gateway
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        body  body      dto.PayBillRequest  true  "Bill Payment Request"
// @Success      201   {object}  dto.RedirectGateway
// @Failure      400   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/payment/pay-bill [post]
func PayUserBill(svcGtr ServiceGetter[paymentp.Service], callbackURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		logPrefix := "PayUserBill handler"

		var req dto.PayBillRequest
		if err := BodyParse(r, &req); err != nil {
			log.Error(logPrefix, zap.Error(err))
			BadRequestError(w, r, err.Error())
			return
		}

		svc := svcGtr(r.Context())
		gateway := paymentd.GatewayType(req.Gateway)
		billID := common.IDFromText(req.BillID)
		userId, ok := r.Context().Value(appjwt.UserIDKey).(string)
		if !ok {
			log.Error(logPrefix, zap.String("error", "failed to get user id from request context"))
			InternalServerError(w, r)
			return
		}
		userID := common.IDFromText(userId)

		resp, err := svc.PayBill(r.Context(), gateway, userID, billID, callbackURL)
		if err != nil {
			log.Error(logPrefix, zap.Error(err))
			switch {
			case errors.Is(err, payment.ErrUnknownGateway):
				Error(w, r, http.StatusBadRequest, payment.ErrUnknownGateway.Error())
			default:
				InternalServerError(w, r)
			}
			return
		}

		if err = WriteJson(w, http.StatusCreated, resp); err != nil {
			log.Error(logPrefix, zap.Error(err))
			InternalServerError(w, r)
		}
	})
}

// PayTotalDebt
//
// @Summary      Pay total debt
// @Description  Initiates payment for all outstanding bills for the authenticated user
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        body  body      dto.PayTotalDebtRequest  true  "Total Debt Payment Request"
// @Success      201   {object}  dto.RedirectGateway
// @Failure      400   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/payment/pay-total-debt [post]
func PayTotalDebt(svcGtr ServiceGetter[paymentp.Service], callbackURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		logPrefix := "PayTotalDebt handler"

		var req dto.PayTotalDebtRequest
		if err := BodyParse(r, &req); err != nil {
			log.Error(logPrefix, zap.Error(err))
			BadRequestError(w, r, err.Error())
			return
		}

		svc := svcGtr(r.Context())
		gateway := paymentd.GatewayType(req.Gateway)
		userId, ok := r.Context().Value(appjwt.UserIDKey).(string)
		if !ok {
			log.Error(logPrefix, zap.String("error", "failed to get user id from request context"))
			InternalServerError(w, r)
			return
		}
		userID := common.IDFromText(userId)

		resp, err := svc.PayTotalDebt(r.Context(), gateway, userID, callbackURL)
		if err != nil {
			log.Error(logPrefix, zap.Error(err))
			switch {
			case errors.Is(err, payment.ErrUnknownGateway):
				Error(w, r, http.StatusBadRequest, payment.ErrUnknownGateway.Error())
			default:
				InternalServerError(w, r)
			}
			return
		}

		if err = WriteJson(w, http.StatusCreated, resp); err != nil {
			log.Error(logPrefix, zap.Error(err))
			InternalServerError(w, r)
		}
	})
}

// CallbackHandler
//
// @Summary      Payment callback
// @Description  Handles payment gateway callback and updates payment status
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        gateway  query    string  true  "Gateway"
// @Param        token    query    string  false "Payment Token"
// @Param        payment-ids query []string false "Payment IDs"
// @Success      200   {object}  dto.PayResponse
// @Failure      400   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/payment/callback [get]
func CallbackHandler(svcGtr ServiceGetter[paymentp.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		logPrefix := "CallbackHandler"

		gatewayStr := r.URL.Query().Get("gateway")
		if gatewayStr == "" {
			BadRequestError(w, r, "missing gateway")
			return
		}
		gateway := paymentd.GatewayType(gatewayStr)

		// Collect all query parameters
		data := map[string][]string{}
		for key, vals := range r.URL.Query() {
			data[key] = vals
		}

		svc := svcGtr(r.Context())
		err := svc.HandleCallback(r.Context(), gateway, data)
		if err != nil {
			log.Error(logPrefix, zap.Error(err))
			Error(w, r, http.StatusBadRequest, err.Error())
			return
		}

		resp := dto.PayResponse{
			Status:  "completed",
			Message: "payment successfully completed",
		}
		if err := WriteJson(w, http.StatusOK, &resp); err != nil {
			log.Error(logPrefix, zap.Error(err))
			InternalServerError(w, r)
		}
	})
}

// SupportedGateways
//
// @Summary      List supported payment gateways
// @Description  Returns a list of supported payment gateways
// @Tags         Payment
// @Produce      json
// @Success      200   {object}  dto.SupportedGatewaysResponse
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/payment/supported-gateways [get]
func SupportedGateways(svcGtr ServiceGetter[paymentp.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())
		logPrefix := "CallbackHandler"
		svc := svcGtr(r.Context())
		resp := &dto.SupportedGatewaysResponse{
			SupportedGateways: svc.SupportedGateways(),
		}
		if err := WriteJson(w, http.StatusOK, resp); err != nil {
			log.Error(logPrefix, zap.Error(err))
			InternalServerError(w, r)
		}
	})
}
