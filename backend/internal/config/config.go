package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                string
	StaticDir           string
	DevAuth             bool
	DiscordClientID     string
	DiscordClientSecret string
}

func Load() (*Config, error) {
	port, err := requireEnv("PORT")
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:                port,
		StaticDir:           getEnv("STATIC_DIR", "../frontend/dist"),
		DevAuth:             os.Getenv("DEV_AUTH") == "true",
		DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	}, nil
}

func (c *Config) HasDiscordCredentials() bool {
	return c.DiscordClientID != "" && c.DiscordClientSecret != ""
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("variável de ambiente obrigatória: %s", key)
	}
	return v, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
