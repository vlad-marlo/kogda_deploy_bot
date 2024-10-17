package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"os"
	"path"
)

type Config struct {
	Telegram TelegramConfig `config:"telegram" toml:"telegram"`
}

func New() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	jsonPath := path.Join(wd, "config.json")
	tomlPath := path.Join(wd, "config.toml")
	yamlPath := path.Join(wd, "config.yaml")
	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewOptionalBackend(tomlPath),
		file.NewOptionalBackend(jsonPath),
		file.NewOptionalBackend(yamlPath),
		flags.NewBackend(),
	)
	config := &Config{
		Telegram: TelegramConfig{
			TelegramTimeoutSeconds: 10,
		},
	}
	err = loader.Load(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
