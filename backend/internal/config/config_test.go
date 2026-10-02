package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFileKeepsExistingEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	content := "# comentário\nTEST_CFG_NEW=\"novo\"\nTEST_CFG_EXISTING=arquivo\ninvalida\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TEST_CFG_EXISTING", "ambiente")
	t.Setenv("TEST_CFG_NEW", "")
	os.Unsetenv("TEST_CFG_NEW")

	loadEnvFile(path)

	if got := os.Getenv("TEST_CFG_NEW"); got != "novo" {
		t.Fatalf("TEST_CFG_NEW = %q, esperado novo", got)
	}
	if got := os.Getenv("TEST_CFG_EXISTING"); got != "ambiente" {
		t.Fatalf("TEST_CFG_EXISTING = %q, esperado ambiente", got)
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("STATIC_DIR", "")
	t.Setenv("DEV_AUTH", "true")

	cfg := Load()

	if cfg.Port != "3000" || cfg.StaticDir != "../frontend/dist" || !cfg.DevAuth {
		t.Fatalf("config inesperada: %+v", cfg)
	}
}
