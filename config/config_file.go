package config

// Support YAML, JSON, TOML files using https://github.com/jinzhu/configor
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

import (
	"fmt"
	"os"

	"github.com/jinzhu/configor"
	"github.com/leeif/mercury/common"
	"github.com/leeif/mercury/server"
	storage "github.com/leeif/mercury/storage/config"
)

const (
	DefaultLogFormat  = "logfmt"
	DefaultLogLevel   = "info"
	DefaultAPIAddress = "127.0.0.1/32"
	DefaultWSAddress  = "0.0.0.0/0"
	DefaultAPIPort    = "6010"
	DefaultWSPort     = "6009"
	DefaultMySQLUser  = "root"
	DefaultMySQLPort  = "3306"
)

type Config struct {
	ConfigFile string
	Log        common.LogConfig
	Server     server.ServerConfig
	Storage    storage.StorageConfig
}

func LoadConfigFile(filePath string, config *Config) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// configure file not exist
		return
	}
	c := struct {
		Log struct {
			Level  string `default:"info"`
			Format string `default:"logfmt"`
		}

		Server struct {
			APIAddress string `default:"127.0.0.1/32"`
			APIPort    string `default:"6009"`
			WSAddress  string `default:"0.0.0.0/0"`
			WSPort     string `default:"6010"`
		}

		MySQL struct {
			Host     string
			Port     string `default:"3306"`
			User     string `default:"root"`
			Password string
		}

		Redis struct {
			Host     string
			Port     string
			User     string
			Password string
		}
	}{}

	fmt.Printf("Loading config file %s...\n", filePath)
	if err := configor.New(&configor.Config{Debug: false}).Load(&c, filePath); err != nil {
		fmt.Printf("Parse config file error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := config.Log.Format.Set(c.Log.Format); err != nil {
		fmt.Printf("Log config error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := config.Log.Level.Set(c.Log.Level); err != nil {
		fmt.Printf("Log config error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := config.Server.APIAddress.Set(c.Server.APIAddress); err != nil {
		fmt.Printf("Server config error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := config.Server.WSAddress.Set(c.Server.WSAddress); err != nil {
		fmt.Printf("Server config error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := config.Server.APIPort.Set(c.Server.APIPort); err != nil {
		fmt.Printf("Server config error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := config.Server.WSPort.Set(c.Server.WSPort); err != nil {
		fmt.Printf("Server config error: %s\n", err.Error())
		os.Exit(1)
	}

	if c.Redis.Host != "" {
		// load redis in the fisrt priority
		config.Storage.RedisConfig.Host = c.Redis.Host
		config.Storage.RedisConfig.Port = c.Redis.Port
		config.Storage.RedisConfig.User = c.Redis.User
		config.Storage.RedisConfig.Password = c.Redis.Password
	} else if c.MySQL.Host != "" {
		config.Storage.MySQLConfig.Host = c.MySQL.Host
		config.Storage.MySQLConfig.Port = c.MySQL.Port
		config.Storage.MySQLConfig.User = c.MySQL.User
		config.Storage.MySQLConfig.Password = c.MySQL.Password
	}
}
