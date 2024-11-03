package observabilityzap

import (
	"context"
	"io"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var _ io.Closer = &tracerProviderCloser{}

type tracerProviderCloser struct {
	tracerProvider *sdktrace.TracerProvider
}

func newTracerProviderCloser(tracerProvider *sdktrace.TracerProvider) *tracerProviderCloser {
	return &tracerProviderCloser{
		tracerProvider: tracerProvider,
	}
}

func (t *tracerProviderCloser) Close() error {
	return t.tracerProvider.Shutdown(context.Background())
}
