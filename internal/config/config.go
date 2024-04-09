package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TgToken          string       `envconfig:"tg_token"`
	MaintainerChatID int64        `envconfig:"maintainer_chat_id"`
	DBType           string       `envconfig:"db_type"`
	Sqlite           SqliteConfig `envconfig:"sqlite"`
}

type SqliteConfig struct {
	File          string `envconfig:"file"`
	MigrationPath string `envconfig:"migration_path"`
}

func InitDotEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("gotenv.Load err: %w", err)
	}
	return nil
}

func Load() (Config, error) {
	var cfg Config
	err := envconfig.Process("app", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("envconfig.Process err: %s", err)
	}
	return cfg, nil
}
