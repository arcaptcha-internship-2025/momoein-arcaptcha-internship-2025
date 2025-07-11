package logger

import "go.uber.org/zap"

type Logger = zap.Logger

type Mode string

const (
	ModeProduction  Mode = "production"
	ModeDevelopment Mode = "development"
)
