package appext

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tracerContainer struct {
	appName string
}

func newTracerContainer(appName string) *tracerContainer {
	return &tracerContainer{
		appName: appName,
	}
}

func (c *tracerContainer) Tracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer(c.appName)
}
