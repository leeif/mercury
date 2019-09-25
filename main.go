package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/leeif/mercury/connection"
	"github.com/leeif/mercury/house"

	"go.uber.org/fx"

	"github.com/leeif/mercury/config"
	mlog "github.com/leeif/mercury/log"
	"github.com/leeif/mercury/server"
	"github.com/leeif/mercury/storage"
)

var VERSION string

func main() {
	app := fx.New(
		fx.Provide(
			func() string {
				return VERSION
			},
			config.NewConfig,
			mlog.NewLogger,
			storage.NewStore,
			house.NewHouse,
			connection.NewPool,
		),
		fx.Invoke(server.Serve),
		fx.NopLogger,
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
