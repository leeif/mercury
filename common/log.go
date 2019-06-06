package common

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	env string
)

type AllowedLevel struct {
	s string
	o level.Option
}

func (l *AllowedLevel) String() string {
	return l.s
}

func (l *AllowedLevel) Set(s string) error {
	switch s {
	case "debug":
		l.o = level.AllowDebug()
	case "info":
		l.o = level.AllowInfo()
	case "warn":
		l.o = level.AllowWarn()
	case "error":
		l.o = level.AllowError()
	default:
		return errors.Errorf("unrecognized log level %q", s)
	}
	l.s = s
	return nil
}

type AllowedFormat struct {
	s string
}

func (f *AllowedFormat) Set(s string) error {
	switch s {
	case "logfmt", "json":
		f.s = s
	default:
		return errors.Errorf("unrecognized log format %q", s)
	}
	return nil
}

func (f *AllowedFormat) String() string {
	return f.s
}

type LogConfig struct {
	Level  *AllowedLevel
	Format *AllowedFormat
}

func SetLogFlag(a *kingpin.Application, config *LogConfig) {
	config.Level = &AllowedLevel{}
	a.Flag("log.level", "[debug, info, warn, error]").
		Default("info").SetValue(config.Level)

	config.Format = &AllowedFormat{}
	a.Flag("log.format", "[logfmt, json]").
		Default("logfmt").SetValue(config.Format)
}

func NewLogger(config *LogConfig) log.Logger {
	var l log.Logger
	if config.Format.s == "logfmt" {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	} else {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	}
	l = level.NewFilter(l, config.Level.o)
	l = log.With(l, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return l
}
