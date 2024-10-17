package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"os"
	"path"
	"time"
)

type Config struct {
	Telegram TelegramConfig `config:"telegram" toml:"telegram"`
	Postgres PostgresConfig `config:"postgres" toml:"postgres"`
	App      AppConfig
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
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}
	config := &Config{
		Telegram: TelegramConfig{
			TelegramTimeoutSeconds: 10,
		},
		Postgres: PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
		},
		App: AppConfig{
			Location: location,
		},
	}
	err = loader.Load(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
