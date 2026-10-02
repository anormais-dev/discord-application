package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("STATIC_DIR", "")
	t.Setenv("DEV_AUTH", "true")

	cfg := Load()

	if cfg.Port != "3000" || cfg.StaticDir != "../frontend/dist" || !cfg.DevAuth {
		t.Fatalf("config inesperada: %+v", cfg)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "4000")
	t.Setenv("STATIC_DIR", "/static")
	t.Setenv("DEV_AUTH", "")
	t.Setenv("DISCORD_CLIENT_ID", "id")
	t.Setenv("DISCORD_CLIENT_SECRET", "secret")

	cfg := Load()

	if cfg.Port != "4000" || cfg.StaticDir != "/static" || cfg.DevAuth || !cfg.HasDiscordCredentials() {
		t.Fatalf("config inesperada: %+v", cfg)
	}
}
