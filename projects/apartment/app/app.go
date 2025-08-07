package app

import (
	"context"
	"database/sql"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/config"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment"
	apartmentPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/apartment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill"
	billPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment"
	paymentd "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	paymentp "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/user"
	userPort "github.com/arcaptcha-internship-2025/momoein-apartment/internal/user/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/email"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/paygw"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/storage"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/minio"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/postgres"
	"go.uber.org/zap"
)

type app struct {
	cfg              config.Config
	logger           *logger.Logger
	db               *sql.DB
	userService      userPort.Service
	apartmentService apartmentPort.Service
	apartmentMail    apartmentPort.EmailSender
	billService      billPort.Service
	paymentService   paymentp.Service
	paymentGateways  map[paymentd.GatewayType]paymentp.Gateway
}

func MustNew(ctx context.Context, cfg config.Config) App {
	app, err := New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return app
}

func New(ctx context.Context, cfg config.Config) (App, error) {
	app := &app{
		cfg:    cfg,
		logger: appctx.Logger(ctx),
	}
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
	app.db = db
	if err = checkMinio(cfg.Minio); err != nil {
		return nil, err
	}
	if err = app.setupPaymentGateways(); err != nil {
		return nil, err
	}
	return app, nil
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
			a.apartmentMailService(),
		)
	}
	return a.apartmentService
}

func (a *app) apartmentMailService() apartmentPort.EmailSender {
	if a.apartmentMail == nil {
		a.apartmentMail = email.NewApartmentEmail(a.cfg.Smaila.Endpoint)
	}
	return a.apartmentMail
}

func checkMinio(cfg config.MinioConfig) error {
	return minio.Ping(
		cfg.Endpoint,
		cfg.AccessKey,
		cfg.SecretKey,
		false,
		3*time.Second)
}

func (a *app) billObjectStorage() (billPort.ObjectStorage, error) {
	c, err := minio.NewClient(
		a.cfg.Minio.Endpoint,
		a.cfg.Minio.AccessKey,
		a.cfg.Minio.SecretKey)
	if err != nil {
		return nil, err
	}
	bos, err := storage.NewBillObjectStorage(c)
	if err != nil {
		return nil, err
	}
	return bos, err
}

func (a *app) BillService() billPort.Service {
	bos, err := a.billObjectStorage()
	if err != nil {
		a.logger.Error("", zap.Error(err))
	}

	if a.billService == nil {
		a.billService = bill.NewService(
			storage.NewBillRepo(a.db),
			bos,
		)
	}
	return a.billService
}

func (a *app) setupPaymentGateways() error {
	gateways := make(map[paymentd.GatewayType]paymentp.Gateway)
	mockGateway, err := paygw.NewMockGateway(a.cfg.BaseURL)
	if err != nil {
		return err
	}
	gateways[paymentd.MockGateway] = mockGateway
	a.paymentGateways = gateways
	return nil
}

func (a *app) PaymentService() paymentp.Service {
	if a.paymentService != nil {
		return a.paymentService
	}
	repo := storage.NewPaymentRepo(a.db)
	gateways := a.paymentGateways
	a.paymentService = payment.NewService(repo, gateways)
	return a.paymentService
}
