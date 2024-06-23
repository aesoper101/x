package appext

import (
	"context"
	"github.com/aesoper101/x/app"
	"github.com/aesoper101/x/configext"
	"github.com/aesoper101/x/internal/verbose"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"path/filepath"
	"time"
)

// NameContainer is a container for named applications.
//
// Application name foo-bar translates to environment variable prefix FOO_BAR_, which is
// used for the various functions that NameContainer provides.
type NameContainer interface {
	// AppName is the application name.
	//
	// The name must be in [a-zA-Z0-9-_].
	AppName() string
	// ConfigDirPath is the config directory path for the named application.
	//
	// First checks for $APP_NAME_CONFIG_DIR.
	// If this is not set, uses app.ConfigDirPath()/app-name.
	// Unnormalized.
	ConfigDirPath() string

	// CacheDirPath is the cache directory path for the named application.
	//
	// First checks for $APP_NAME_CACHE_DIR.
	// If this is not set, uses app.CacheDirPath()/app-name.
	// Unnormalized.
	CacheDirPath() string
	// DataDirPath is the data directory path for the named application.
	//
	// First checks for $APP_NAME_DATA_DIR.
	// If this is not set, uses app.DataDirPath()/app-name.
	// Unnormalized.
	DataDirPath() string
	// Port is the port to use for serving.
	//
	// First checks for $APP_NAME_PORT.
	// If this is not set, checks for $PORT.
	// If this is not set, returns 0, which means no port is known.
	// Returns error on parse.
	Port() (uint16, error)
}

// NewNameContainer returns a new NameContainer.
//
// The name must be in [a-zA-Z0-9-_].
func NewNameContainer(envContainer app.EnvContainer, appName string) (NameContainer, error) {
	return newNameContainer(envContainer, appName)
}

// LoggerContainer provides a *zap.Logger.
type LoggerContainer interface {
	Logger() *slog.Logger
}

// NewLoggerContainer returns a new LoggerContainer.
func NewLoggerContainer(logger *slog.Logger) LoggerContainer {
	return newLoggerContainer(logger)
}

// TracerContainer provides a trace.Tracer based on the application name.
type TracerContainer interface {
	Tracer() trace.Tracer
}

// NewTracerContainer returns a new TracerContainer for the application name.
func NewTracerContainer(appName string) TracerContainer {
	return newTracerContainer(appName)
}

// VerboseContainer provides a verbose.Printer.
type VerboseContainer interface {
	// VerboseEnabled returns true if verbose mode is enabled.
	VerboseEnabled() bool
	// VerbosePrinter returns a verbose.Printer to use for verbose printing.
	VerbosePrinter() verbose.Printer
}

// NewVerboseContainer returns a new VerboseContainer.
func NewVerboseContainer(verbosePrinter verbose.Printer) VerboseContainer {
	return newVerboseContainer(verbosePrinter)
}

// Container contains not just the base app container, but all extended containers.
type Container interface {
	app.Container
	NameContainer
	LoggerContainer
	TracerContainer
	VerboseContainer
}

// NewContainer returns a new Container.
func NewContainer(
	baseContainer app.Container,
	appName string,
	logger *slog.Logger,
	verbosePrinter verbose.Printer,
) (Container, error) {
	return newContainer(
		baseContainer,
		appName,
		logger,
		verbosePrinter,
	)
}

// Interceptor intercepts and adapts the request or response of run functions.
type Interceptor func(func(context.Context, Container) error) func(context.Context, Container) error

// SubCommandBuilder builds run functions for sub-commands.
type SubCommandBuilder interface {
	NewRunFunc(func(context.Context, Container) error) func(context.Context, app.Container) error
}

// Builder builds run functions for both top-level commands and sub-commands.
type Builder interface {
	BindRoot(flagSet *pflag.FlagSet)
	SubCommandBuilder
}

// BuilderOption is an option for a new Builder
type BuilderOption func(*builder)

// BuilderWithTimeout returns a new BuilderOption that adds a timeout flag and the default timeout.
func BuilderWithTimeout(defaultTimeout time.Duration) BuilderOption {
	return func(builder *builder) {
		builder.defaultTimeout = defaultTimeout
	}
}

// BuilderWithTracing enables zap tracing for the builder.
func BuilderWithTracing() BuilderOption {
	return func(builder *builder) {
		builder.tracing = true
	}
}

// BuilderWithInterceptor adds the given interceptor for all run functions.
func BuilderWithInterceptor(interceptor Interceptor) BuilderOption {
	return func(builder *builder) {
		builder.interceptors = append(builder.interceptors, interceptor)
	}
}

// BuilderWithDefaultLogLevel adds the given default log level.
func BuilderWithDefaultLogLevel(defaultLogLevel slog.Level) BuilderOption {
	return func(builder *builder) {
		builder.defaultLogLevel = defaultLogLevel
	}
}

func ReadConfig(ctx context.Context, container NameContainer, value interface{}) error {
	configFilePath := filepath.Clean(container.ConfigDirPath())
	provider, err := configext.New(ctx, configext.WithConfigFiles(configFilePath), configext.EnableEnvLoading(container.AppName()))
	if err != nil {
		return err
	}

	return provider.Unmarshal("", value)
}
