package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user"
	userDomain "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/domain"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrBadRequest     = errors.New("invalid request payload")
	ErrInternalServer = errors.New("internal server error")
)

func getSignUpHandler(service userPort.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req dto.SignUpRequest
		body := r.Body
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		u, err := service.Create(context.Background(),
			UserDTOToDomain(&dto.User{Email: req.Email, Password: req.Password}))
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

		token, err := GenerateJWT([]byte("secret"), u.Email.String())
		if err != nil {
			e := fmt.Errorf("%w: %w", ErrInternalServer, err)
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}

		resp := &dto.SignUpResponse{
			AccessToken: token,
		}
		bResp, _ := json.Marshal(resp)
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
		})
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(bResp)
	})
}

func UserDTOToDomain(u *dto.User) *userDomain.User {
	id, err := uuid.FromBytes([]byte(u.ID))
	if err != nil {
		id = userDomain.NilID
	}
	user := &userDomain.User{
		ID:        id,
		Email:     userDomain.Email(u.Email),
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
	user.SetPassword([]byte(u.Password))
	return user
}

func GenerateJWT(secret []byte, email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(secret)
}
