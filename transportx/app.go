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

type appInfo struct {
	id       string
	name     string
	version  string
	metadata map[string]string
}

func NewAppInfo(id, name, version string, metadata map[string]string) AppInfo {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &appInfo{id, name, version, metadata}
}

// ID returns app instance id.
func (app *appInfo) ID() string { return app.id }

// Name returns service name.
func (app *appInfo) Name() string { return app.name }

// Version returns app version.
func (app *appInfo) Version() string { return app.version }

// Metadata returns service metadata.
func (app *appInfo) Metadata() map[string]string { return app.metadata }

type runner struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger *slog.Logger

	appInfo AppInfo

	sigs        []os.Signal
	stopTimeout time.Duration

	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error

	servers []Server

	hasStopped bool
}

func newRunner(ctx context.Context, appInfo AppInfo, logger *slog.Logger, opts ...RunOption) *runner {
	app := &runner{
		ctx:    context.Background(),
		logger: logger,

		appInfo: appInfo,

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

func (app *runner) Run() error {
	sctx := NewContext(app.ctx, app.appInfo)
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
			stopCtx, cancel := context.WithTimeout(NewContext(app.ctx, app.appInfo), app.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(NewContext(app.ctx, app.appInfo))
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

func (app *runner) Stop() error {
	if app.hasStopped {
		app.logger.Warn("app: app already stopped", "id", app.appInfo.ID())
		return nil
	}

	sctx := NewContext(app.ctx, app.appInfo)
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

func Run(ctx context.Context, appInfo AppInfo, logger *slog.Logger, options ...RunOption) error {
	r := newRunner(ctx, appInfo, logger, options...)
	return r.Run()
}
