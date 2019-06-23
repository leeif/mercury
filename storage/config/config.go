package config

type StorageConfig struct {
	MaxUnRead   int
	MySQLConfig *MySQLConfig
	RedisConfig *RedisConfig
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}
