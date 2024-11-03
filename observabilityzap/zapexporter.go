package observabilityzap

import (
	"context"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var _ trace.SpanExporter = &zapExporter{}

type zapExporter struct {
	logger *zap.Logger
}

func newZapExporter(logger *zap.Logger) *zapExporter {
	return &zapExporter{
		logger: logger,
	}
}

// ExportSpans 将一组trace.ReadOnlySpan导出到zap日志中
func (z *zapExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for _, span := range spans {
		if z.shouldLogSpan(span) {
			z.logSpan(span)
		}
	}
	return nil
}

// shouldLogSpan 判断是否应该记录该span
func (z *zapExporter) shouldLogSpan(span trace.ReadOnlySpan) bool {
	if !span.SpanContext().IsSampled() {
		return false
	}
	return z.logger.Check(zap.DebugLevel, span.Name()) != nil
}

// logSpan 将span记录到zap日志中
func (z *zapExporter) logSpan(span trace.ReadOnlySpan) {
	fields := z.createSpanFields(span)
	z.logger.Debug(span.Name(), fields...)
}

// createSpanFields 创建span的zap字段列表
func (z *zapExporter) createSpanFields(span trace.ReadOnlySpan) []zap.Field {
	fields := []zap.Field{
		zap.Duration("duration", span.EndTime().Sub(span.StartTime())),
		zap.String("status", span.Status().Code.String()),
	}
	for _, attribute := range span.Attributes() {
		fields = append(fields, zap.Any(string(attribute.Key), attribute.Value.AsInterface()))
	}
	for _, event := range span.Events() {
		for _, attribute := range event.Attributes {
			// Event attributes seem to have their event name magically prepended to the attribute key.
			// This could overlap with attributes, but we're going to ignore this for now since it's extremely unlikely,
			// and since this is only really for the CLI. Not a good answer.
			fields = append(fields, zap.Any(string(attribute.Key), attribute.Value.AsInterface()))
		}
	}
	return fields
}

func (z *zapExporter) Shutdown(ctx context.Context) error {
	return nil
}
