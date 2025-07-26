package app

import (
	"context"
	"database/sql"

	"github.com/arcaptcha-internship-2025/momoein-apartment/config"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment"
	apartmentPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill"
	billPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/storage"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/minio"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/postgres"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/smtp"
)

type app struct {
	cfg              config.Config
	logger           *logger.Logger
	db               *sql.DB
	userService      userPort.Service
	apartmentService apartmentPort.Service
	apartmentMail    apartmentPort.Email
	billService      billPort.Service
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
		User:    cfg.DB.User,
		Pass:    cfg.DB.Password,
		Host:    cfg.DB.Host,
		Port:    cfg.DB.Port,
		DBName:  cfg.DB.DBName,
		Schema:  cfg.DB.Schema,
		AppName: cfg.DB.AppName,
	}
	db, err := postgres.NewPSQLConn(opt)
	if err != nil {
		return nil, err
	}

	return &app{
		cfg:    cfg,
		db:     db,
		logger: appctx.Logger(ctx),
	}, nil
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) Logger() *logger.Logger {
	return a.logger
}
func (a *app) DB() *sql.DB {
	return a.db
}

func (a *app) UserService(ctx context.Context) userPort.Service {
	if a.userService == nil {
		a.userService = user.NewService(storage.NewUserRepo(a.db))
	}
	return a.userService
}

func (a *app) ApartmentService(ctx context.Context) apartmentPort.Service {
	if a.apartmentService == nil {
		a.apartmentService = apartment.NewService(
			storage.NewApartmentRepo(a.db),
			a.mailService(),
		)
	}
	return a.apartmentService
}

func (a *app) mailService() apartmentPort.Email {
	c := a.Config().SMTP
	if a.apartmentMail == nil {
		a.apartmentMail = smtp.NewSMTPService(c.Host, c.Port, c.From, c.Username, c.Password)
	}
	return a.apartmentMail
}

func (a *app) BillService() billPort.Service {
	c := minio.MustNewClient(
		a.cfg.Minio.Endpoint,
		a.cfg.Minio.AccessKey,
		a.cfg.Minio.SecretKey,
	)
	if a.billService == nil {
		a.billService = bill.NewService(
			storage.NewBillRepo(a.db),
			storage.MustNewBillObjectStorage(c),
		)
	}
	return a.billService
}
