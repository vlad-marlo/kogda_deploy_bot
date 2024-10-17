package config

import "fmt"

type PostgresConfig struct {
	Host     string `config:"postgres-host" toml:"host" json:"host"`
	Port     int    `config:"postgres-port" toml:"port" json:"port"`
	User     string `config:"postgres-user" toml:"user" json:"user"`
	Password string `config:"postgres-password" toml:"password" json:"password"`
	Database string `config:"postgres-database" toml:"database" json:"database"`
}

func (cfg PostgresConfig) URI() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
}
