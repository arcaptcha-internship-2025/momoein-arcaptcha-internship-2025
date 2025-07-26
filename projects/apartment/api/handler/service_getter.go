package handler

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
	apartmentPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	billPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
)

type ServiceGetter[T any] func(context.Context) T

func UserServiceGetter(app app.App) ServiceGetter[userPort.Service] {
	return func(ctx context.Context) userPort.Service {
		return app.UserService(ctx)
	}
}

func ApartmentServiceGetter(a app.App) ServiceGetter[apartmentPort.Service] {
	return func(ctx context.Context) apartmentPort.Service {
		return a.ApartmentService(ctx)
	}
}

func BillServiceGetter(a app.App) ServiceGetter[billPort.Service] {
	return func(ctx context.Context) billPort.Service {
		return a.BillService()
	}
}
