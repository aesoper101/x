package transportx

import (
	"context"
	"os"
	"time"
)

type RunOption func(*runner)

func WithServers(server ...Server) RunOption {
	return func(app *runner) {
		app.servers = append(app.servers, server...)
	}
}

func WithSignal(signal ...os.Signal) RunOption {
	return func(app *runner) {
		app.sigs = append(app.sigs, signal...)
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
