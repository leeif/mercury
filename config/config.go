package config

import (
	"os"
	"path/filepath"

	"github.com/leeif/kiper"
)

// Config file
// Support YAML, JSON, TOML files
// Example of a toml configure file
// mc.cnf.toml
/*
[log]
level="info"

[server]
apiAddress="127.0.0.1/32"
wsAddress="0.0.0.0/0"
port="6010"

[mysql]
host="127.0.0.1"
user="root"
*/

type Config struct {
	Log     LogConfig
	Server  ServerConfig
	Storage StorageConfig
}

func NewConfig(version string) (*Config, error) {
	k := kiper.NewKiper(filepath.Base(os.Args[0]), "mercury server")
	k.SetConfigFileFlag("config.file", "config file", "./config.toml")
	c := Config{
		Log:     newLogConfig(),
		Server:  newServerConfig(),
		Storage: newStorageConfig(),
	}
	k.Kingpin.Version(version)
	k.Kingpin.HelpFlag.Short('h')
	err := k.Parse(&c, os.Args[1:])
	if err != nil {
		return nil, err
	}
	return &c, nil
}
