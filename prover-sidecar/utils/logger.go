package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateLogger(verbose bool) *zap.Logger {
	logLevel := zapcore.InfoLevel
	if verbose {
		logLevel = zapcore.DebugLevel
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	loggerConfig.Encoding = "console"
	loggerConfig.Level = zap.NewAtomicLevelAt(logLevel)

	// Create the logger from the core
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
