package config

import "os"

type Config struct {
	Port                string
	StaticDir           string
	DevAuth             bool
	DiscordClientID     string
	DiscordClientSecret string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "3000"),
		StaticDir:           getEnv("STATIC_DIR", "../frontend/dist"),
		DevAuth:             os.Getenv("DEV_AUTH") == "true",
		DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	}
}

func (c *Config) HasDiscordCredentials() bool {
	return c.DiscordClientID != "" && c.DiscordClientSecret != ""
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
