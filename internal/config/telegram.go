package config

import "time"

type TelegramConfig struct {
	TelegramSecretKey      string `config:"telegram-secret-key,required" toml:"secret_key"`
	TelegramTimeoutSeconds int    `config:"telegram-timeout-seconds,short=t" toml:"timeout_seconds"`
}

func (cfg TelegramConfig) Timeout() time.Duration {
	return time.Duration(cfg.TelegramTimeoutSeconds) * time.Second
}
