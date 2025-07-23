package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env    string    `yaml:"env" env-required:"true"`
	DbCfg  DBConfig  `yaml:"db_config" env-required:"true"`
	ApiCfg APIConfig `yaml:"api_config" env-required:"true"`
}

type DBConfig struct {
	Addr                string        `yaml:"address" env-required:"true"`
	Username            string        `yaml:"username" env-required:"true"`
	Password            string        `yaml:"password" env-required:"true"`
	DbName              string        `yaml:"db_name" env-required:"true"`
	SslMode             string        `yaml:"ssl_mode" env-required:"true"`
	PoolMaxConns        int           `yaml:"pool_max_conns" env-required:"true"`
	PoolMaxConnLifetime time.Duration `yaml:"pool_max_conn_lifetime" env-required:"true"`
	PoolMaxConnIdleTime time.Duration `yaml:"pool_max_conn_idletime" env-required:"true"`
}

type APIConfig struct {
	Addr        string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
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
