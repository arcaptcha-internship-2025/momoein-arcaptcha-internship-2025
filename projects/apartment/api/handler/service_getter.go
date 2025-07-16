package handler

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
)

type ServiceGetter[T any] func(context.Context) T

func UserServiceGetter(app app.App) ServiceGetter[userPort.Service] {
	return func(ctx context.Context) userPort.Service {
		return app.UserService(ctx)
	}
}
