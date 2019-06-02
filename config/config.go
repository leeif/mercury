package conf

import (
	"io/ioutil"

	"github.com/leeif/mercury/common"
	"github.com/leeif/mercury/server"
	sConfig "github.com/leeif/mercury/storage/config"
)

type Config struct {
	ConfigFile    string
	LogConfig     common.LogConfig
	ServerConfig  server.ServerConfig
	StorageConfig sConfig.StorageConfig
}

func LoadConfigFile(filepath string, config Config) {
	_, err := ioutil.ReadFile(filepath)
	if err != nil {

	}
}
