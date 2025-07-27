package main

import (
	"context"
	"flag"
	"os"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler"
	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
	"github.com/arcaptcha-internship-2025/momoein-apartment/config"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	"go.uber.org/zap"
)

var envfile = flag.String("env-file", ".env", "environment file path.")

func main() {
	flag.Parse()

	if v := os.Getenv("ENV_FILE"); len(v) > 0 {
		*envfile = v
	}
	cfg := config.MustReadEnv(*envfile)

	appLogger := logger.NewZapLogger(logger.ModeProduction)
	if cfg.AppMode == config.Development {
		appLogger = logger.NewZapLogger(logger.ModeDevelopment)
	}
	defer appLogger.Sync()

	ctx := appctx.New(context.Background(), appctx.WithLogger(appLogger))

	appContainer := app.MustNew(ctx, cfg)
	appLogger.Info("Application started")
	appLogger.Fatal("", zap.Error(handler.Run(appContainer)))
}
