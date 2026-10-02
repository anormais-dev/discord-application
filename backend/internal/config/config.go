package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port                string
	StaticDir           string
	DevAuth             bool
	DiscordClientID     string
	DiscordClientSecret string
}

func Load() *Config {
	loadEnvFile("../.env")

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

// loadEnvFile não sobrescreve o que já está no ambiente.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, strings.Trim(strings.TrimSpace(value), `"'`))
		}
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
