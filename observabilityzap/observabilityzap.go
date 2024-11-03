package observabilityzap

import (
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Start creates a Zap logging exporter for Opentelemetry traces and returns
// the exporter. The exporter implements io.Closer for clean-up.
func Start(logger *zap.Logger) (trace.TracerProvider, io.Closer) {
	exporter := newZapExporter(logger)
	tracerProviderOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	}
	tracerProvider := sdktrace.NewTracerProvider(tracerProviderOptions...)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tracerProvider, newTracerProviderCloser(tracerProvider)
}
