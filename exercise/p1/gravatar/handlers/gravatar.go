package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/mail"
)

type GravatarResponse struct {
	Ok          bool   `json:"ok"`
	GravatarUrl string `json:"gravatar_url"`
}

type ErrorResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

var (
	ErrNoEmail = ErrorResponse{
		Ok:      false,
		Message: "No email address provided",
	}
	ErrInvalidEmail = ErrorResponse{
		Ok:      false,
		Message: "Invalid email address",
	}
)

func HandleGravatarRequest(w http.ResponseWriter, r *http.Request) {
	emailString := r.URL.Query().Get("email")
	if len(emailString) <= 0 {
		WriteResponseJSON(w, http.StatusBadRequest, &ErrNoEmail)
		return
	}
	if _, err := mail.ParseAddress(emailString); err != nil {
		WriteResponseJSON(w, http.StatusBadRequest, &ErrInvalidEmail)
		return
	}

	hash := md5.Sum([]byte(emailString))
	hashString := hex.EncodeToString(hash[:])

	WriteResponseJSON(w, http.StatusOK, &GravatarResponse{
		Ok:          true,
		GravatarUrl: "https://www.gravatar.com/avatar/" + hashString,
	})
}

func WriteResponseJSON(w http.ResponseWriter, code int, body any) {
	msg, _ := json.Marshal(body)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(msg)
}
