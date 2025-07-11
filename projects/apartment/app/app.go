package app

import (
	"context"
	"database/sql"

	"github.com/arcaptcha-internship-2025/momoein-apartment/config"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/postgres"
)

type app struct {
	cfg    config.Config
	logger *logger.Logger
	db     *sql.DB
}

func MustNew(ctx context.Context, cfg config.Config) App {
	app, err := New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return app
}

func New(ctx context.Context, cfg config.Config) (App, error) {
	opt := postgres.DBConnOptions{
		User:   cfg.DB.User,
		Pass:   cfg.DB.Password,
		Host:   cfg.DB.Host,
		Port:   cfg.DB.Port,
		DBName: cfg.DB.DBName,
		Schema: cfg.DB.Schema,
	}
	db, err := postgres.NewPSQLConn(opt)
	if err != nil {
		return nil, err
	}

	return &app{
		db: db,
		logger: appctx.Logger(ctx),
	}, nil
}

func (a *app) Config(ctx context.Context) config.Config {
	return a.cfg
}

func (a *app) Logger(ctx context.Context) *logger.Logger {
	return a.logger
}
func (a *app) DB(ctx context.Context) *sql.DB {
	return a.db
}
