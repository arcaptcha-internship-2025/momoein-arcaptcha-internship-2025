package app

import (
	"context"
	"database/sql"

	"github.com/arcaptcha-internship-2025/momoein-apartment/config"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
)

type App interface {
	Config() config.Config
	Logger() *logger.Logger
	DB() *sql.DB
	UserService(ctx context.Context) port.Service
}
