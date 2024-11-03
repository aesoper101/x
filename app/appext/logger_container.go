package appext

import (
	"go.uber.org/zap"
)

type loggerContainer struct {
	logger *zap.Logger
}

func newLoggerContainer(logger *zap.Logger) *loggerContainer {
	return &loggerContainer{
		logger: logger,
	}
}

func (c *loggerContainer) Logger() *zap.Logger {
	return c.logger
}
