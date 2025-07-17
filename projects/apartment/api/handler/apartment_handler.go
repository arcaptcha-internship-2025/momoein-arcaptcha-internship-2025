package handler

import (
	"net/http"

	apartmentPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
)

func AddApartment(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func InviteApartmentMember(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func ApartmentMembers(svcGetter ServiceGetter[apartmentPort.Service]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

