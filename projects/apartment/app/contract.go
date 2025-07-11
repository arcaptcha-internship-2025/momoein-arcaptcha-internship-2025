package app

import (
	"context"
	"database/sql"

	"github.com/arcaptcha-internship-2025/momoein-apartment/config"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
)

type App interface {
	Config(ctx context.Context) config.Config
	Logger(ctx context.Context) *logger.Logger
	DB(ctx context.Context) *sql.DB
}
