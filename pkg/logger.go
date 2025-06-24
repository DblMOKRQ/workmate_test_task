package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New создает новый экземпляр логгера
func NewLogger() (*zap.Logger, error) {

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	baseLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return baseLogger, nil
}
