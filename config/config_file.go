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
	"reflect"

	"github.com/jinzhu/configor"
	"github.com/leeif/mercury/common"
	"github.com/leeif/mercury/server"
	storage "github.com/leeif/mercury/storage/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
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
		// return
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

		Storage struct {
			MaxUnRead int `default:"100"`
		}

		MySQL struct {
			Host     string
			Port     string `default:"3306"`
			User     string `default:"root"`
			Password string
		}

		Redis struct {
			Host     string
			Port     string `default:"6379"`
			Password string
			DB       int `default:"0"`
		}
	}{}

	fmt.Printf("Loading config file %s...\n", filePath)
	if err := configor.New(&configor.Config{Debug: false}).Load(&c, filePath); err != nil {
		fmt.Printf("Parse config file error: %s\n", err.Error())
		os.Exit(1)
	}

	setConfigKingpin(config.Log.Format, c.Log.Format)
	setConfigKingpin(config.Log.Level, c.Log.Level)

	setConfigKingpin(config.Server.APIAddress, c.Server.APIAddress)
	setConfigKingpin(config.Server.APIPort, c.Server.APIPort)
	setConfigKingpin(config.Server.WSAddress, c.Server.WSAddress)
	setConfigKingpin(config.Server.WSPort, c.Server.WSPort)

	setConfig(&config.Storage, c.Storage)
	setConfig(config.Storage.MySQLConfig, c.MySQL)
	setConfig(config.Storage.RedisConfig, c.Redis)
}

func setConfigKingpin(target kingpin.Value, value string) {
	if target.String() == "" {
		err := target.Set(value)
		configError(err)
	}
}

func setConfig(target interface{}, value interface{}) {
	targetValue := reflect.Indirect(reflect.ValueOf(target))
	valueType := reflect.TypeOf(value)
	valueValue := reflect.ValueOf(value)
	for i := 0; i < valueType.NumField(); i++ {
		name := valueType.Field(i).Name
		tvField := targetValue.FieldByName(name)

		if !tvField.IsValid() || !tvField.CanSet() {
			continue
		}
		vField := valueValue.FieldByName(name)
		switch tvField.Kind() {
		case reflect.String:
			if targetValue.FieldByName(name).String() == "" && vField.Kind() == reflect.String {
				targetValue.FieldByName(name).SetString(vField.String())
			}
		case reflect.Int:
			if targetValue.FieldByName(name).Int() == -1 && vField.Kind() == reflect.Int {
				targetValue.FieldByName(name).SetInt(vField.Int())
			}
		}

	}
}

func configError(err error) {
	if err != nil {
		fmt.Printf("Server config error: %s\n", err.Error())
		os.Exit(1)
	}
}
