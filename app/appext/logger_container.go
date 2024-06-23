package appext

import "log/slog"

type loggerContainer struct {
	logger *slog.Logger
}

func newLoggerContainer(logger *slog.Logger) *loggerContainer {
	return &loggerContainer{
		logger: logger,
	}
}

func (c *loggerContainer) Logger() *slog.Logger {
	return c.logger
}
