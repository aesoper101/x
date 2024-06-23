package appext

import (
	"github.com/aesoper101/x/app"
	"github.com/aesoper101/x/internal/verbose"
	"log/slog"
)

type container struct {
	app.Container
	NameContainer
	LoggerContainer
	TracerContainer
	VerboseContainer
}

func newContainer(
	baseContainer app.Container,
	appName string,
	logger *slog.Logger,
	verbosePrinter verbose.Printer,
) (*container, error) {
	nameContainer, err := newNameContainer(baseContainer, appName)
	if err != nil {
		return nil, err
	}
	return &container{
		Container:        baseContainer,
		NameContainer:    nameContainer,
		LoggerContainer:  newLoggerContainer(logger),
		TracerContainer:  newTracerContainer(appName),
		VerboseContainer: newVerboseContainer(verbosePrinter),
	}, nil
}
