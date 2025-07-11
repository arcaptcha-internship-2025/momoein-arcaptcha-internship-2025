package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(mode Mode) *Logger {
	logLevel := zap.InfoLevel
	if mode == ModeDevelopment {
		logLevel = zap.DebugLevel
	}

	level := zap.NewAtomicLevelAt(logLevel)
	consoleEncoderCfg := getConsoleEncoderConfig()
	fileEncoderCfg := getFileEncoderConfig()

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderCfg)
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderCfg)

	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    20, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	opts := []zap.Option{zap.AddCaller()}
	if mode == ModeDevelopment {
		opts = append(opts, zap.Development())
	}

	return zap.New(core, opts...)
}

func getConsoleEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.TimeKey = "timestamp"
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg
}

func getFileEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.TimeKey = "timestamp"
	return cfg
}
