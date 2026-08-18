// Package app configures and runs application.
package app

import (
	"context"

	"github.com/keep/sunny/config"
	"github.com/keep/sunny/internal/modules"
	"github.com/keep/sunny/internal/observability"
	"github.com/keep/sunny/pkg/httpserver"
	"github.com/keep/sunny/pkg/logger"
	"go.uber.org/fx"
)

// New builds the application's dependency graph without starting it.
func New(cfg config.Config) *fx.App {
	return fx.New(
		fx.Supply(
			fx.Annotate(
				cfg,
				fx.As(new(config.Config)),
			),
		),

		infrastructureModule,
		transportModule,
		// api
		modules.APIModule,
		// observability
		observability.Module,
		// start http server
		fx.Invoke(startHTTPServer),
	)
}

// Run starts the application and blocks until shutdown.
func Run(cfg config.Config) {
	New(cfg).Run()
}

func startHTTPServer(
	lc fx.Lifecycle,
	logger logger.Logger,
	httpServer httpserver.Server,
) {
	lc.Append(
		fx.Hook{
			// start
			OnStart: func(ctx context.Context) error {
				logger.Info("http server - Starting HTTP Server...")

				if err := httpServer.Start(); err != nil {
					return err
				}

				logger.Info("http server - Listening on ", httpServer.GetAddress())
				return nil
			},
			// stop
			OnStop: func(ctx context.Context) error {
				logger.Info("http server - Stopping HTTP Server...")

				err := httpServer.Shutdown()
				if err != nil {
					return err
				}

				logger.Info("http server - Server - Shutting down")
				return nil
			},
		},
	)
}
