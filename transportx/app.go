package transportx

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
}

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger *slog.Logger

	// app info
	id       string
	name     string
	version  string
	metadata map[string]string

	sigs        []os.Signal
	stopTimeout time.Duration

	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error

	servers []Server

	hasStopped bool
}

func New(opts ...Option) *App {
	app := &App{
		ctx:    context.Background(),
		logger: slog.Default(),

		sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		stopTimeout: time.Second * 10,
	}
	for _, opt := range opts {
		opt(app)
	}

	ctx, cancel := context.WithCancel(app.ctx)
	app.ctx = ctx
	app.cancel = cancel

	return app
}

// ID returns app instance id.
func (app *App) ID() string { return app.id }

// Name returns service name.
func (app *App) Name() string { return app.name }

// Version returns app version.
func (app *App) Version() string { return app.version }

// Metadata returns service metadata.
func (app *App) Metadata() map[string]string { return app.metadata }

func (app *App) Run() error {
	sctx := NewContext(app.ctx, app)
	eg, ctx := errgroup.WithContext(sctx)
	wg := sync.WaitGroup{}

	for _, fn := range app.beforeStart {
		if err := fn(sctx); err != nil {
			app.logger.Error("app: before start error", "err", err)
			return err
		}
	}

	for _, srv := range app.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx, cancel := context.WithTimeout(NewContext(app.ctx, app), app.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(NewContext(app.ctx, app))
		})
	}
	wg.Wait()

	for _, fn := range app.afterStart {
		if err := fn(sctx); err != nil {
			app.logger.Error("app: after start error", "err", err)
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, app.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return app.Stop()
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		app.logger.Error("app: run error", "err", err)
		return err
	}

	for _, fn := range app.afterStop {
		if err := fn(sctx); err != nil {
			app.logger.Error("app: after stop error", "err", err)
			return err
		}
	}

	return nil
}

func (app *App) Stop() error {
	if app.hasStopped {
		app.logger.Warn("app: app already stopped", "id", app.id)
		return nil
	}

	sctx := NewContext(app.ctx, app)
	for _, fn := range app.beforeStop {
		if err := fn(sctx); err != nil {
			app.logger.Error("app: before stop error", "err", err)
			return err
		}
	}

	app.hasStopped = true
	if app.cancel != nil {
		app.cancel()
	}

	return nil
}
