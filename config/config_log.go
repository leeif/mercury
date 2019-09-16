package config

import (
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type LogConfig struct {
	Level  *AllowedLevel  `kiper_value:"name:level;help:[debug, info, warn, error];default:info"`
	Format *AllowedFormat `kiper_value:"name:format;help:[logfmt, json];default:logfmt"`
}

type AllowedLevel struct {
	S string
	O level.Option
}

func (l *AllowedLevel) String() string {
	return l.S
}

func (l *AllowedLevel) Set(s string) error {
	switch s {
	case "debug":
		l.O = level.AllowDebug()
	case "info":
		l.O = level.AllowInfo()
	case "warn":
		l.O = level.AllowWarn()
	case "error":
		l.O = level.AllowError()
	default:
		return errors.Errorf("unrecognized log level %q", s)
	}
	l.S = s
	return nil
}

type AllowedFormat struct {
	S string
}

func (f *AllowedFormat) Set(s string) error {
	switch s {
	case "logfmt", "json":
		f.S = s
	default:
		return errors.Errorf("unrecognized log format %q", s)
	}
	return nil
}

func (f *AllowedFormat) String() string {
	return f.S
}

func newLogConfig() LogConfig {
	return LogConfig{
		Level:  &AllowedLevel{},
		Format: &AllowedFormat{},
	}
}
