package log

import (
	"os"

	"github.com/leeif/mercury/config"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func NewLogger(config *config.Config) log.Logger {
	var l log.Logger
	cfg := config.Log
	if cfg.Format.S == "logfmt" {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	} else {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	}
	l = level.NewFilter(l, cfg.Level.O)
	l = log.With(l, "ts", log.DefaultTimestampUTC)
	if cfg.Level.S == "debug" {
		l = log.With(l, "caller", log.DefaultCaller)
	}
	return l
}
