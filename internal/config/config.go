package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	DB       DBConfig       `yaml:"database"`
	Telegram TelegramConfig `yaml:"telegram"`
	Root     RootConfig     `yaml:"root"`
}

type AppConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Secret       string `yaml:"secret"`
	IsProduction bool   `yaml:"is_production"`
	RunSeeder    bool   `yaml:"run_seeder"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type TelegramConfig struct {
	BotToken   string `yaml:"bot_token"`
	BotEnabled bool   `yaml:"bot_enabled"`
	GroupID    int    `yaml:"group_id"`
	MiniAppURL string `yaml:"mini_app_url"`
}

type RootConfig struct {
	TID       int64     `yaml:"tid"`
	UUID      uuid.UUID `yaml:"uuid"`
	FirstName string    `yaml:"first_name"`
	LastName  string    `yaml:"last_name"`
	Username  string    `yaml:"username"`
	Email     string    `yaml:"email"`
	Phone     string    `yaml:"phone"`
	Password  string    `yaml:"password"`
}

func New(path string) *Config {
	cfg, err := newConfig(path)
	if err != nil {
		slog.Warn("Failed to load config", slog.Any("err", err))
		os.Exit(1)
	}
	return cfg
}

func newConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}
	var cfg Config
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.yaml: %w", err)
	}
	return &cfg, nil
}
