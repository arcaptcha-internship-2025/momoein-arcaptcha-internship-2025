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

var configPath = flag.String("config", "config.json", "configuration file path, it must be json")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_FILE"); len(v) > 0 {
		*configPath = v
	}
	cfg := config.MustReadJson(*configPath)

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
