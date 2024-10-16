package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default: "local"`
	DBPath   string        `yaml:"db_path" env-requred: "true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-requred: "true"`
	GRPC     GRPCConfig    `yaml: "grpc"`
}

type GRPCConfig struct {
	Port   int           `yaml:"port"`
	Timout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		panic("config file does not exist, err:" + err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config, err: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
