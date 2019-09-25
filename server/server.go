package server

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/leeif/mercury/house"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/mercury/config"
	"go.uber.org/fx"
)

func Serve(lc fx.Lifecycle, config *config.Config, logger log.Logger, house *house.House) {
	cfg := config.Server
	logger = log.With(logger, "component", "server")

	apiRouter := newAPIRouter(cfg, house, logger)

	apiAddress := ":" + cfg.API.Port.String()
	apiListener, err := net.Listen("tcp", apiAddress)

	if err != nil {
		level.Error(logger).Log("error", err.Error())
		os.Exit(1)
	}

	apiServer := &http.Server{
		Handler: apiRouter,
	}

	wsRouter := newWSRouter(cfg, house, logger)

	wsAddress := ":" + cfg.WS.Port.String()
	wsListener, err := net.Listen("tcp", wsAddress)

	if err != nil {
		level.Error(logger).Log("error", err.Error())
		os.Exit(1)
	}

	wsServer := &http.Server{
		Handler: wsRouter,
	}

	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 30 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go func() {
				level.Info(logger).Log("msg", "WebSocket server is listening at "+wsAddress)
				wsServer.Serve(wsListener)
			}()
			go func() {
				level.Info(logger).Log("msg", "API server is listening at "+apiAddress)
				apiServer.Serve(apiListener)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			level.Info(logger).Log("Stopping Pluto server")
			if err := wsServer.Shutdown(ctx); err != nil {
				return err
			}
			if err := apiServer.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		},
	})
}
