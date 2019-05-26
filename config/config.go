package conf

import (
	"io/ioutil"
	"github.com/leeif/mercury/common"
	"github.com/leeif/mercury/server"
	"github.com/leeif/mercury/storage"
)

type Config struct {
	ConfigFile    string
	LogConfig     common.LogConfig
	ServerConfig  server.ServerConfig
	StorageConfig storage.StorageConfig
}

func LoadConfigFile(filepath string, config Config) {
	_, err := ioutil.ReadFile(filepath)
	if err != nil {

	}
}
