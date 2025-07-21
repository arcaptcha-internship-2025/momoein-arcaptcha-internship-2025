package handler

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment"
	apartmentPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	appjwt "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/jwt"
	"go.uber.org/zap"
)

func AddApartment(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		var req = new(dto.Apartment)
		if err := BodyParse(r, req); err != nil {
			Error(w, r, http.StatusBadRequest, "BadRequest", err.Error())
			return
		}

		svc := svcGetter(r.Context())

		userID, ok := r.Context().Value(appjwt.UserIDKey).(string)
		if !ok {
			log.Warn("failed to get user id from request context")
			Error(w, r, http.StatusUnauthorized, "Unauthorized")
			return
		}
		req.AdminID = userID

		a, err := svc.Create(r.Context(), dto.ApartmentDTOToDomain(req))
		if err != nil {
			log.Error("AddApartment", zap.Error(err))
			Error(w, r, http.StatusInternalServerError, "InternalServerError", err.Error())
			return
		}

		if err = WriteJson(w, http.StatusCreated, a); err != nil {
			log.Error("failed to write response", zap.Error(err))
			InternalServerError(w, r)
		}
	})
}

func InviteApartmentMember(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		var req = new(dto.InviteUserToApartmentRequest)
		if err := BodyParse(r, req); err != nil {
			log.Warn("body parse", zap.Error(err))
			Error(w, r, http.StatusBadRequest, "BadRequest", err.Error())
			return
		}
		svc := svcGetter(r.Context())

		id, ok := r.Context().Value(appjwt.UserIDKey).(string)
		if !ok {
			log.Error("request context", zap.String("Error", "failed to get user id from request context"))
			Error(w, r, http.StatusInternalServerError, "InternalServerError")
			return
		}
		adminId := common.NilID
		if err := adminId.UnmarshalText([]byte(id)); err != nil {
			log.Error("user id", zap.Error(err))
			Error(w, r, http.StatusInternalServerError, "InternalServerError")
			return
		}

		member, err := svc.InviteMember(r.Context(), adminId, req.ApartmentID, common.Email(req.UserEmail))
		if err != nil {
			log.Error("invite member", zap.Error(err))
			Error(w, r, http.StatusInternalServerError, "InternalServerError")
			return
		}

		if err = WriteJson(w, http.StatusOK, member); err != nil {
			log.Error("WriteJson response", zap.Error(err))
			InternalServerError(w, r)
		}
	})
}

const InviteTokenKey string = "token"

func AcceptApartmentInvite(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		if !r.URL.Query().Has(InviteTokenKey) {
			log.Info("invite token not exists")
			Error(w, r, http.StatusBadRequest, "token not exists")
			return
		}
		token := r.URL.Query().Get(InviteTokenKey)

		svc := svcGetter(r.Context())
		err := svc.AcceptInvite(r.Context(), token)
		if err != nil {
			switch {
			case errors.Is(err, apartment.ErrInvalidToken):
				Error(w, r, http.StatusBadRequest)
			case errors.Is(err, apartment.ErrUnregisteredUser): // ??
				redirectURL := "/sign-up?return_to=" + url.QueryEscape(r.URL.String())
				http.Redirect(w, r, redirectURL, http.StatusFound)
			default:
				log.Error("accept invite", zap.Error(err))
				Error(w, r, http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})
}

func ApartmentMembers(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, r, http.StatusNotImplemented, "Not Implemented")
	})
}
