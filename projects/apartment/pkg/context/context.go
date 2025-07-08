package appctx

import (
	"context"

	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
)

type CtxKey string

var defaultLogger *logger.Logger

func init() {
	defaultLogger = logger.NewZapLogger(logger.ModeProduction)
}

type appContext struct {
	context.Context
	logger *logger.Logger
}

type AppContextOpt func(*appContext) *appContext

func New(parent context.Context, opts ...AppContextOpt) context.Context {
	ctx := &appContext{Context: parent}
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	return ctx
}

func WithLogger(logger *logger.Logger) AppContextOpt {
	return func(ctx *appContext) *appContext {
		ctx.logger = logger
		return ctx
	}
}

func SetLogger(ctx context.Context, logger *logger.Logger) {
	if appCtx, ok := ctx.(*appContext); ok {
		appCtx.logger = logger
	}
}

func Logger(ctx context.Context) *logger.Logger {
	appCtx, ok := ctx.(*appContext)
	if !ok || appCtx.logger == nil {
		return defaultLogger
	}
	return appCtx.logger
}

func WithValue(key CtxKey, val any) AppContextOpt {
	return func(ac *appContext) *appContext {
		ac.Context = context.WithValue(ac.Context, key, val)
		return ac
	}
}

func SetValue(ctx context.Context, key CtxKey, val any) {
	if appCtx, ok := ctx.(*appContext); ok {
		appCtx.Context = context.WithValue(appCtx.Context, key, val)
	}
}

func IsAppContext(ctx context.Context) bool {
	_, ok := ctx.(*appContext)
	return ok
}
