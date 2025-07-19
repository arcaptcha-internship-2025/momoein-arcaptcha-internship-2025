package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/config"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	appjwt "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrBadRequest     = errors.New("invalid request payload")
	ErrInternalServer = errors.New("internal server error")
)

func getSignUpHandler(svcGetter ServiceGetter[userPort.Service], cfg config.AuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := appctx.Logger(r.Context())

		var req dto.SignUpRequest
		if err := BodyParse(r, &req); err != nil {
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

		authResp, err := GenerateAuthResponse(cfg, u.ID.String(), u.Email.String())
		if err != nil {
			log.Error("failed to generate jwt token", zap.Error(err))
			Error(w, r, http.StatusInternalServerError, "InternalServerError", err.Error())
			return
		}

		SetTokenCookie(w, cfg, authResp.AccessToken, authResp.RefreshToken)

		if err = WriteJson(w, http.StatusCreated, authResp); err != nil {
			log.Error("failed to write response", zap.Error(err))
			Error(w, r, http.StatusInternalServerError, "InternalServerError", err.Error())
		}
	})
}

func getSignInHandler(svcGetter ServiceGetter[userPort.Service], cfg config.AuthConfig) http.Handler {
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
			log.Warn("password comparison failed", zap.Error(err))
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		authResp, err := GenerateAuthResponse(cfg, u.ID.String(), u.Email.String())
		if err != nil {
			log.Error("failed to generate jwt token", zap.Error(err))
			Error(w, r, http.StatusInternalServerError, "InternalServerError", err.Error())
			return
		}

		SetTokenCookie(w, cfg, authResp.AccessToken, authResp.RefreshToken)

		if err = WriteJson(w, http.StatusCreated, authResp); err != nil {
			log.Error("failed to write response", zap.Error(err))
			Error(w, r, http.StatusInternalServerError, "InternalServerError", err.Error())
		}
	})
}

func GenerateAuthResponse(cfg config.AuthConfig, id, email string) (*dto.AuthResponse, error) {
	accessToken, err := createJWTToken(cfg.JWTSecret, id, email, cfg.AccessExpiry)
	if err != nil {
		return nil, err
	}
	refreshToken, err := createJWTToken(cfg.JWTSecret, id, email, cfg.RefreshExpiry)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func createJWTToken(secret, userId, userEmail string, exp int64) (string, error) {
	return appjwt.CreateToken([]byte(secret), &appjwt.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * time.Duration(exp))},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		UserID:    userId,
		UserEMail: userEmail,
	})
}

func SetTokenCookie(
	w http.ResponseWriter,
	cfg config.AuthConfig,
	accessToken, refreshToken string,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * time.Duration(cfg.AccessExpiry)),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Minute * time.Duration(cfg.RefreshExpiry)),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})
}
