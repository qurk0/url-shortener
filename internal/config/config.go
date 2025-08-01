package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env       string        `yaml:"env" env-required:"true"`
	ApiCfg    APIConfig     `yaml:"api_config" env-required:"true"`
	MainDBCfg PGSQLConfig   `yaml:"main_db_config" env-required:"true"`
	CacheCfg  RedisConfig   `yaml:"cache_config" env-required:"true"`
	Clients   ClientsConfig `yaml:"clients"`
	AppSecret []byte        `yaml:"app_secret" env-required:"true"`
}

type PGSQLConfig struct {
	Host                string        `yaml:"host" env-required:"true"`
	Port                string        `yaml:"port" env-required:"true"`
	Username            string        `yaml:"username" env-required:"true"`
	Password            string        `yaml:"password" env-required:"true"`
	DbName              string        `yaml:"db_name" env-required:"true"`
	SslMode             string        `yaml:"ssl_mode" env-required:"true"`
	PoolMaxConns        int           `yaml:"pool_max_conns" env-required:"true"`
	PoolMaxConnLifetime time.Duration `yaml:"pool_max_conn_lifetime" env-required:"true"`
	PoolMaxConnIdleTime time.Duration `yaml:"pool_max_conn_idletime" env-required:"true"`
}

type RedisConfig struct {
	Addr     string        `yaml:"addr" env-required:"true"`
	Password string        `yaml:"password" env-required:"true"`
	DB       int           `yaml:"db" env-default:"0"`
	TTL      time.Duration `yaml:"ttl" env-default:"12h"`
}

type APIConfig struct {
	Addr        string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type Client struct {
	Address    string        `yaml:"address"`
	Timeout    time.Duration `yaml:"timeout"`
	RetryCount int           `yaml:"retry_count"`
}

type ClientsConfig struct {
	Auth Client `yaml:"auth"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	path := os.Getenv("CONFIG_FILE_PATH")
	if path == "" {
		panic("CONFIG_FILE_PATH not set")
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist")
	}

	var cfg Config
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
