package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrBadRequest     = errors.New("invalid request payload")
	ErrInternalServer = errors.New("internal server error")
)

func getSignUpHandler(svcGetter ServiceGetter[userPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		var req dto.SignUpRequest
		body := r.Body
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		service := svcGetter(r.Context())
		u, err := service.Create(r.Context(),
			dto.UserDTOToDomain(&dto.User{Email: req.Email, Password: req.Password}))
		if err != nil {
			switch {
			case errors.Is(err, user.ErrUserOnValidate):
				e := fmt.Errorf("%w: %w", ErrBadRequest, err)
				http.Error(w, e.Error(), http.StatusBadRequest)
			default:
				e := fmt.Errorf("%w: %w", ErrInternalServer, err)
				http.Error(w, e.Error(), http.StatusInternalServerError)
			}
			return
		}

		token, err := GenerateJWT(secret, u.Email.String())
		if err != nil {
			e := fmt.Errorf("%w: %w", ErrInternalServer, err)
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(&dto.AuthResponse{AccessToken: token})
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			http.Error(w, ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

		SetContentTypeJson(w)
		SetTokenCookie(w, token)
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	})
}

func getSignInHandler(svcGetter ServiceGetter[userPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		var req dto.SignInRequest
		if err := BodyParse(r, &req); err != nil {
			http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		service := svcGetter(r.Context())
		u, err := service.Get(r.Context(), &domain.UserFilter{Email: common.Email(req.Email)})
		if err != nil {
			switch {
			case errors.Is(err, user.ErrUserNotFound):
				http.Error(w, user.ErrUserNotFound.Error(), http.StatusNotFound)
			default:
				e := fmt.Errorf("%w: %w", ErrInternalServer, err)
				http.Error(w, e.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err := u.ComparePassword([]byte(req.Password)); err != nil {
			log.Warn("compare password", zap.String("password", req.Password), zap.Error(err))
			http.Error(w, "incorrect password", http.StatusUnauthorized)
			return
		}

		token, err := GenerateJWT(secret, req.Email)
		if err != nil {
			log.Error("failed to generate jwt token", zap.Error(err))
			http.Error(w, ErrInternalServer.Error(), http.StatusInternalServerError)
			return 
		}

		resp, err := json.Marshal(&dto.AuthResponse{AccessToken: token})
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			http.Error(w, ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

		SetContentTypeJson(w)
		SetTokenCookie(w, token)
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	})
}

var secret = []byte("this-is-very-secret")

func GenerateJWT(secret []byte, email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(secret)
}

func BodyParse(r *http.Request, dest any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&dest)
}

func SetContentTypeJson(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}

func SetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})
}
