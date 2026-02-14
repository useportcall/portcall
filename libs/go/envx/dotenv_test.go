package envx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPrefersDotEnv(t *testing.T) {
	tmp := t.TempDir()
	prev, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".env"), []byte("EMAIL_PROVIDER=resend\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".env.example"), []byte("EMAIL_PROVIDER=local\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_ENV", "development")
	if err := os.Unsetenv("EMAIL_PROVIDER"); err != nil {
		t.Fatal(err)
	}
	Load()
	if got := os.Getenv("EMAIL_PROVIDER"); got != "resend" {
		t.Fatalf("expected EMAIL_PROVIDER=resend, got %q", got)
	}
}

func TestLoadFallsBackToExample(t *testing.T) {
	tmp := t.TempDir()
	prev, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".env.example"), []byte("EMAIL_FROM=hello@useportcall.com\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_ENV", "development")
	if err := os.Unsetenv("EMAIL_FROM"); err != nil {
		t.Fatal(err)
	}
	Load()
	if got := os.Getenv("EMAIL_FROM"); got != "hello@useportcall.com" {
		t.Fatalf("expected EMAIL_FROM from .env.example, got %q", got)
	}
}

func TestLoadPrefersDotEnvs(t *testing.T) {
	tmp := t.TempDir()
	prev, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".envs"), []byte("DISCORD_WEBHOOK_URL_SIGNUP=from_envs\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".env"), []byte("DISCORD_WEBHOOK_URL_SIGNUP=from_env\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".env.example"), []byte("DISCORD_WEBHOOK_URL_SIGNUP=from_example\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_ENV", "development")
	if err := os.Unsetenv("DISCORD_WEBHOOK_URL_SIGNUP"); err != nil {
		t.Fatal(err)
	}
	Load()
	if got := os.Getenv("DISCORD_WEBHOOK_URL_SIGNUP"); got != "from_envs" {
		t.Fatalf("expected DISCORD_WEBHOOK_URL_SIGNUP from .envs, got %q", got)
	}
}
