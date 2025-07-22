package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/dto"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"go.uber.org/zap"
)

func BodyParse(r *http.Request, dest any) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(dest)
}

func SetContentTypeJson(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}

func WriteJson(w http.ResponseWriter, code int, body any) error {
	SetContentTypeJson(w)
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}

func Error(w http.ResponseWriter, r *http.Request, code int, msg ...string) {
	log := appctx.Logger(r.Context())

	body := dto.Error{
		Code:    code,
		Message: strings.Join(msg, ": "),
	}

	resp, err := json.Marshal(body)
	if err != nil {
		log.Error("failed to marshal error response", zap.Error(err))
		http.Error(w, `{"code":500,"message":"Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	SetContentTypeJson(w)
	w.WriteHeader(code)
	_, _ = w.Write(resp)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusInternalServerError, "Internal Server Error")
}

func BadRequestError(w http.ResponseWriter, r *http.Request, msg ...string) {
	m := append([]string{"Bad Request"}, msg...)
	Error(w, r, http.StatusBadRequest, m...)
}
