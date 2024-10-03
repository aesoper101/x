package transportx

//go:generate mockgen -package=transportx -destination=mock_server_test.go . Server

import (
	"context"
	"errors"
	"github.com/aesoper101/x/interrupt"
	"github.com/aesoper101/x/uuidutil"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/url"
	"sync"
	"time"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Runner interface {
	Run() error
	Stop() error
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

	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	stopTimeout time.Duration

	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error

	servers []Server

	hasStopped bool
}

func newRunner(opts ...RunOption) *runner {
	app := &runner{
		ctx:    context.Background(),
		logger: slog.Default(),

		stopTimeout: time.Second * 10,
	}

	if id, err := uuidutil.New(); err == nil {
		app.id = id.String()
	}

	for _, opt := range opts {
		opt(app)
	}

	app.appInfo = &appInfo{
		id:       app.id,
		name:     app.name,
		version:  app.version,
		metadata: app.metadata,
	}

	ctx, cancel := interrupt.WithCancel(app.ctx)
	app.ctx, app.cancel = ctx, cancel

	return app
}

func (app *runner) Run() (err error) {
	sctx := NewContext(app.ctx, app.appInfo)
	eg, ctx := errgroup.WithContext(sctx)
	wg := sync.WaitGroup{}

	for _, fn := range app.beforeStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}

	for _, srv := range app.servers {
		srv := srv
		eg.Go(
			func() error {
				<-ctx.Done() // wait for stop signal
				stopCtx, cancel := context.WithTimeout(NewContext(app.ctx, app.appInfo), app.stopTimeout)
				defer cancel()
				return srv.Stop(stopCtx)
			},
		)
		wg.Add(1)
		eg.Go(
			func() error {
				wg.Done() // here is to ensure server start has begun running before register, so defer is not needed
				startCtx := NewContext(app.ctx, app.appInfo)
				return srv.Start(startCtx)
			},
		)
	}
	wg.Wait()

	for _, fn := range app.afterStart {
		if err = fn(sctx); err != nil {
			return err
		}
	}

	eg.Go(
		func() error {
			select {
			case <-app.ctx.Done():
				return app.Stop()
			}
		},
	)
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	err = nil
	for _, fn := range app.afterStop {
		err = fn(sctx)
	}
	return err
}

func (app *runner) Stop() (err error) {
	if app.hasStopped {
		app.logger.Warn("runner has stopped")
		return nil
	}

	sctx := NewContext(app.ctx, app.appInfo)
	for _, fn := range app.beforeStop {
		err = fn(sctx)
	}

	if app.cancel != nil {
		app.cancel()
	}

	app.logger.Info("runner stopped")

	app.hasStopped = true
	return err
}

func Run(opts ...RunOption) error {
	app := newRunner(opts...)
	return app.Run()
}

func NewRunner(opts ...RunOption) Runner {
	return newRunner(opts...)
}
