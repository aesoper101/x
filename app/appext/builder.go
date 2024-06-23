package appext

import (
	"context"
	"fmt"
	"github.com/aesoper101/x/app"
	"github.com/aesoper101/x/internal/verbose"
	"github.com/spf13/pflag"
	"log/slog"
	"time"
)

type builder struct {
	appName string

	verbose   bool
	debug     bool
	noWarn    bool
	logFormat string

	profile           bool
	profilePath       string
	profileLoops      int
	profileType       string
	profileAllowError bool

	parallelism int

	timeout time.Duration

	defaultTimeout time.Duration

	tracing bool

	// 0 is InfoLevel in slog
	defaultLogLevel slog.Level
	interceptors    []Interceptor
}

func newBuilder(appName string, options ...BuilderOption) *builder {
	builder := &builder{
		appName:         appName,
		defaultLogLevel: slog.LevelInfo,
	}
	for _, option := range options {
		option(builder)
	}
	return builder
}

func (b *builder) BindRoot(flagSet *pflag.FlagSet) {
	flagSet.BoolVarP(&b.verbose, "verbose", "v", false, "Turn on verbose mode")
	flagSet.BoolVar(&b.debug, "debug", false, "Turn on debug logging")
	flagSet.StringVar(&b.logFormat, "log-format", "color", "The log format [text,json]")
	if b.defaultTimeout > 0 {
		flagSet.DurationVar(&b.timeout, "timeout", b.defaultTimeout, `The duration until timing out, setting it to zero means no timeout`)
	}

	// We do not officially support this flag, this is for testing, where we need warnings turned off.
	flagSet.BoolVar(&b.noWarn, "no-warn", false, "Turn off warn logging")
	_ = flagSet.MarkHidden("no-warn")
}

func (b *builder) NewRunFunc(
	f func(context.Context, Container) error,
) func(context.Context, app.Container) error {
	interceptor := chainInterceptors(b.interceptors...)
	return func(ctx context.Context, appContainer app.Container) error {
		if interceptor != nil {
			return b.run(ctx, appContainer, interceptor(f))
		}
		return b.run(ctx, appContainer, f)
	}
}

func (b *builder) run(
	ctx context.Context,
	appContainer app.Container,
	f func(context.Context, Container) error,
) (retErr error) {
	logLevel, err := getLogLevel(b.defaultLogLevel, b.debug, b.noWarn)
	if err != nil {
		return err
	}
	var logHandler slog.Handler
	if b.logFormat == "" {
		logHandler = slog.NewJSONHandler(appContainer.Stderr(), &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		logHandler = slog.NewTextHandler(appContainer.Stderr(), &slog.HandlerOptions{
			Level: logLevel,
		})
	}
	logger := slog.New(logHandler)

	verbosePrinter := verbose.NewPrinterForFlagValue(appContainer.Stderr(), b.appName, b.verbose)
	container, err := newContainer(appContainer, b.appName, logger, verbosePrinter)
	if err != nil {
		return err
	}

	var cancel context.CancelFunc
	if !b.profile && b.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, b.timeout)
		defer cancel()
	}

	if b.tracing {
		// TODO: We should probably have a way to configure the tracer.
	}

	return f(ctx, container)
}

func getLogLevel(defaultLogLevel slog.Level, debugFlag bool, noWarnFlag bool) (slog.Level, error) {
	if debugFlag && noWarnFlag {
		return slog.LevelInfo, fmt.Errorf("cannot set both --debug and --no-warn")
	}
	if noWarnFlag {
		return slog.LevelError, nil
	}
	if debugFlag {
		return slog.LevelDebug, nil
	}
	return defaultLogLevel, nil
}

// chainInterceptors consolidates the given interceptors into one.
// The interceptors are applied in the order they are declared.
func chainInterceptors(interceptors ...Interceptor) Interceptor {
	if len(interceptors) == 0 {
		return nil
	}
	filtered := make([]Interceptor, 0, len(interceptors))
	for _, interceptor := range interceptors {
		if interceptor != nil {
			filtered = append(filtered, interceptor)
		}
	}
	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		first := filtered[0]
		return func(next func(context.Context, Container) error) func(context.Context, Container) error {
			for i := len(filtered) - 1; i > 0; i-- {
				next = filtered[i](next)
			}
			return first(next)
		}
	}
}
