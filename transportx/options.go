package transportx

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type Option func(*App)

func WithServer(server ...Server) Option {
	return func(app *App) {
		app.servers = append(app.servers, server...)
	}
}

func WithContext(ctx context.Context) Option {
	return func(app *App) {
		app.ctx = ctx
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(app *App) {
		app.logger = logger
	}
}

func WithName(name string) Option {
	return func(app *App) {
		app.name = name
	}
}

func WithVersion(version string) Option {
	return func(app *App) {
		app.version = version
	}
}

func WithID(id string) Option {
	return func(app *App) {
		app.id = id
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(app *App) {
		app.metadata = metadata
	}
}

func WithSignal(signal ...os.Signal) Option {
	return func(app *App) {
		app.sigs = append(app.sigs, signal...)
	}
}

func WithStopTimeout(stopTimeout time.Duration) Option {
	return func(app *App) {
		app.stopTimeout = stopTimeout
	}
}

func BeforeStart(fn ...func(ctx context.Context) error) Option {
	return func(app *App) {
		app.beforeStart = append(app.beforeStart, fn...)
	}
}

func AfterStart(fn ...func(ctx context.Context) error) Option {
	return func(app *App) {
		app.afterStart = append(app.afterStart, fn...)
	}
}

func BeforeStop(fn ...func(ctx context.Context) error) Option {
	return func(app *App) {
		app.beforeStop = append(app.beforeStop, fn...)
	}
}

func AfterStop(fn ...func(ctx context.Context) error) Option {
	return func(app *App) {
		app.afterStop = append(app.afterStop, fn...)
	}
}
