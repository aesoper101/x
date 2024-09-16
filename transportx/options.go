package transportx

import (
	"context"
	"log/slog"
	"net/url"
	"time"
)

type RunOption func(*runner)

func WithServers(server ...Server) RunOption {
	return func(app *runner) {
		app.servers = append(app.servers, server...)
	}
}

func WithStopTimeout(stopTimeout time.Duration) RunOption {
	return func(app *runner) {
		app.stopTimeout = stopTimeout
	}
}

func BeforeStart(fn ...func(ctx context.Context) error) RunOption {
	return func(app *runner) {
		app.beforeStart = append(app.beforeStart, fn...)
	}
}

func AfterStart(fn ...func(ctx context.Context) error) RunOption {
	return func(app *runner) {
		app.afterStart = append(app.afterStart, fn...)
	}
}

func BeforeStop(fn ...func(ctx context.Context) error) RunOption {
	return func(app *runner) {
		app.beforeStop = append(app.beforeStop, fn...)
	}
}

func AfterStop(fn ...func(ctx context.Context) error) RunOption {
	return func(app *runner) {
		app.afterStop = append(app.afterStop, fn...)
	}
}

func WithLogger(logger *slog.Logger) RunOption {
	return func(app *runner) {
		app.logger = logger
	}
}

func WithContext(ctx context.Context) RunOption {
	return func(app *runner) {
		app.ctx = ctx
	}
}

func ID(id string) RunOption {
	return func(app *runner) {
		app.id = id
	}
}

func Name(name string) RunOption {
	return func(app *runner) {
		app.name = name
	}
}

func Version(version string) RunOption {
	return func(app *runner) {
		app.version = version
	}
}

func Metadata(metadata map[string]string) RunOption {
	return func(app *runner) {
		app.metadata = metadata
	}
}

func Endpoints(endpoints []*url.URL) RunOption {
	return func(app *runner) {
		app.endpoints = endpoints
	}
}
