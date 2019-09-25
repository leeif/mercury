package config

import "errors"

type StorageConfig struct {
	Type        *StorageType `kiper_value:"name:type;help:[memory, mysql, redis];default:memory"`
	MySQLConfig MySQLConfig  `kiper_config:"name:mysql"`
	RedisConfig RedisConfig  `kiper_config:"name:redis"`
}

type StorageType struct {
	Type string
}

func (st *StorageType) Set(t string) error {
	switch t {
	case "memory":
	case "mysql":
	case "redis":
		st.Type = t
	default:
		return errors.New("Unsupport storage type")
	}
	return nil
}

func (st *StorageType) String() string {
	return st.Type
}

type MySQLConfig struct {
	Host     string `kiper_value:"name:host;help:mysql host;default:127.0.0.1"`
	Port     string `kiper_value:"name:port;help:mysql port;default:3306"`
	User     string `kiper_value:"name:user;help:mysql user;default:root"`
	Password string `kiper_value:"name:password;help:mysql password"`
}

type RedisConfig struct {
	Host     string `kiper_value:"name:host;help:redis host;default:127.0.0.1"`
	Port     string `kiper_value:"name:port;help:redis port;default:6379"`
	Password string `kiper_value:"name:host;help:redis password"`
	DB       int    `kiper_value:"name:host;help:redis db;default:0"`
}

func newStorageConfig() StorageConfig {
	return StorageConfig{
		Type:        &StorageType{},
		MySQLConfig: MySQLConfig{},
		RedisConfig: RedisConfig{},
	}
}
