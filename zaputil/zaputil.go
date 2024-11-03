package zaputil

import (
	"fmt"
	cond "github.com/aesoper101/x/condition"
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Format string

const (
	FormatText  Format = "text"
	FormatColor Format = "color"
	FormatJSON  Format = "json"
)

func (f Format) String() string {
	return string(f)
}

func (f Format) IsValid() bool {
	switch f {
	case FormatText, FormatColor, FormatJSON:
		return true
	default:
		return false
	}
}

// NewLogger returns a new Logger.
func NewLogger(opts ...Option) *zap.Logger {
	options := newOptions(opts...)
	if options.encoder == nil {
		encoder, _ := getZapEncoder(options.format)
		options.encoder = cond.Ternary(encoder != nil, encoder, NewJSONEncoder())
	}

	return zap.New(
		zapcore.NewCore(
			options.encoder,
			zapcore.Lock(zapcore.AddSync(options.writer)),
			zap.NewAtomicLevelAt(options.level),
		),
		options.zapOptions...,
	)
}

// NewLoggerForFlagValues returns a new Logger for the given level and format strings.
//
// The level can be [debug,info,warn,error,panic,fatal]. The default is info.
// The format can be [text,color,json]. The default is color.
func NewLoggerForFlagValues(writer io.Writer, levelString string, format string) (*zap.Logger, error) {
	level, err := getZapLevel(levelString)
	if err != nil {
		return nil, err
	}
	return NewLogger(
		WithWriter(writer),
		WithLevel(level),
		WithFormat(Format(format)),
	), nil
}

// NewTextEncoder returns a new text Encoder.
func NewTextEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(textEncoderConfig)
}

// NewColortextEncoder returns a new colortext Encoder.
func NewColortextEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(colortextEncoderConfig)
}

// NewJSONEncoder returns a new JSON encoder.
func NewJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(jsonEncoderConfig)
}

func getZapLevel(level string) (zapcore.Level, error) {
	level = strings.TrimSpace(strings.ToLower(level))
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info", "":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	default:
		return 0, fmt.Errorf("unknown log level [debug,info,warn,error,panic,fatal]: %q", level)
	}
}

func getZapEncoder(format Format) (zapcore.Encoder, error) {
	f := strings.TrimSpace(strings.ToLower(format.String()))
	switch Format(f) {
	case FormatText:
		return NewTextEncoder(), nil
	case FormatColor, "":
		return NewColortextEncoder(), nil
	case FormatJSON:
		return NewJSONEncoder(), nil
	default:
		return nil, fmt.Errorf("unknown log format [text,color,json]: %q", format)
	}
}
