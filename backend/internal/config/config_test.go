package config

import "testing"

func TestLoadRequiresPort(t *testing.T) {
	t.Setenv("PORT", "")

	if _, err := Load(); err == nil {
		t.Fatalf("Load deveria falhar sem PORT")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "3000")
	t.Setenv("STATIC_DIR", "")
	t.Setenv("DEV_AUTH", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.StaticDir != "../frontend/dist" || !cfg.DevAuth {
		t.Fatalf("config inesperada: %+v", cfg)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "4000")
	t.Setenv("STATIC_DIR", "/static")
	t.Setenv("DEV_AUTH", "")
	t.Setenv("DISCORD_CLIENT_ID", "id")
	t.Setenv("DISCORD_CLIENT_SECRET", "secret")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "4000" || cfg.StaticDir != "/static" || cfg.DevAuth || !cfg.HasDiscordCredentials() {
		t.Fatalf("config inesperada: %+v", cfg)
	}
}
