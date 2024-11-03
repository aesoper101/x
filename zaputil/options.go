package zaputil

import (
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type options struct {
	level   zapcore.Level
	writer  io.Writer
	encoder zapcore.Encoder
	format  Format
}

type Option func(*options)

func newOptions(opts ...Option) *options {
	o := &options{
		level:  zapcore.InfoLevel,
		writer: os.Stderr,
		format: FormatText,
	}
	o.apply(opts...)
	return o
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithLevel(level zapcore.Level) Option {
	return func(o *options) {
		if _, err := zapcore.ParseLevel(level.String()); err != nil {
			return
		}
		o.level = level
	}
}

func WithWriter(writer io.Writer) Option {
	return func(o *options) {
		if writer != nil {
			o.writer = writer
		}
	}
}

// WithEncoder sets the zapcore.Encoder for the logger or encoder.
// Note that the encoder cannot be used with WithFormat at the same time.
func WithEncoder(encoder zapcore.Encoder) Option {
	return func(o *options) {
		if encoder != nil {
			o.encoder = encoder
		}
	}
}

// WithFormat sets the format for the logger or encoder.
// Note that the format cannot be used with WithEncoder at the same time. if both are set, WithEncoder will prevail.
func WithFormat(format Format) Option {
	return func(o *options) {
		if format.IsValid() {
			o.format = format
		}
	}
}
