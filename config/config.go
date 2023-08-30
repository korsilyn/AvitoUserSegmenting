package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
	"time"
)

type (
	Config struct {
		App `yaml:"app"`
		HTTP `yaml:"http"`
		Log `yaml:"log"`
		PG `yaml:"postgres"`
	}

	App struct {
		Name string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"rue" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		MaxPoolSize int `env-required:"true" yaml:"max_pool_size" env:"PG_MAX_POOL_SIZE"`
		URL string `env-required:"true" env:"PG_URL"`
	}
}

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("Error reading CFG file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error updating env: %w", err)
	}

	return cfg, nil
}
