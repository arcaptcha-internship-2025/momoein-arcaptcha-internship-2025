package handler

import (
	"errors"
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment"
	apartmentPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	appjwt "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/jwt"
	"go.uber.org/zap"
)

// AddApartment
//
// @Summary      Create a new apartment
// @Description  Adds a new apartment and assigns the current user as admin
// @Tags         Apartment
// @Accept       json
// @Produce      json
// @Param        body  body      dto.Apartment  true  "Apartment Info"
// @Success      201   {object}  dto.Apartment
// @Failure      400   {object}  dto.Error
// @Failure      401   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/apartment [post]
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

// InviteApartmentMember
//
// @Summary      Invite user to apartment
// @Description  Sends an invitation to a user to join an apartment
// @Tags         Apartment
// @Accept       json
// @Produce      json
// @Param        body  body      dto.InviteUserToApartmentRequest  true  "Invite Request"
// @Success      200   {object}  dto.InviteUserToApartmentResponse
// @Failure      400   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/apartment/invite [post]
func InviteApartmentMember(svcGetter ServiceGetter[apartmentPort.Service], acceptURL string) http.Handler {
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

		member, err := svc.InviteMember(r.Context(), adminId, req.ApartmentID, common.Email(req.UserEmail), acceptURL)
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

// AcceptApartmentInvite
//
// @Summary      Accept apartment invitation
// @Description  Accepts an invitation to join an apartment using a token
// @Tags         Apartment
// @Accept       json
// @Produce      json
// @Param        token  query    string  true  "Invitation Token"
// @Success      202   {string}  string  "Accepted"
// @Failure      400   {object}  dto.Error
// @Failure      401   {object}  dto.Error
// @Failure      500   {object}  dto.Error
// @Router       /api/v1/apartment/invite/accept [get]
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
			case errors.Is(err, apartment.ErrUnregisteredUser):
				Error(w, r, http.StatusUnauthorized, "unregistered user")
			default:
				log.Error("accept invite", zap.Error(err))
				Error(w, r, http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})
}

// ApartmentMembers
//
// @Summary      List apartment members
// @Description  Returns a list of users in the apartment (not implemented)
// @Tags         Apartment
// @Produce      json
// @Failure      501   {object}  dto.Error
// @Router       /api/v1/apartment/members [get]
func ApartmentMembers(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, r, http.StatusNotImplemented, "Not Implemented")
	})
}
