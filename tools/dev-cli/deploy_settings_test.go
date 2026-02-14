package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDeploySettingsMissingFile(t *testing.T) {
	t.Parallel()
	settings, err := loadDeploySettings(t.TempDir())
	if err != nil {
		t.Fatalf("loadDeploySettings returned error: %v", err)
	}
	if settings.SkipMigration {
		t.Fatal("expected default SkipMigration=false")
	}
}

func TestSaveAndLoadDeploySettings(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	want := deploySettings{SkipMigration: true}
	if err := saveDeploySettings(dir, want); err != nil {
		t.Fatalf("saveDeploySettings returned error: %v", err)
	}
	got, err := loadDeploySettings(dir)
	if err != nil {
		t.Fatalf("loadDeploySettings returned error: %v", err)
	}
	if got != want {
		t.Fatalf("settings mismatch: got %+v want %+v", got, want)
	}
}

func TestLoadDeploySettingsInvalidJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, ".dev-cli.deploy.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("failed to write invalid settings file: %v", err)
	}
	if _, err := loadDeploySettings(dir); err == nil {
		t.Fatal("expected parse error for invalid deploy settings json")
	}
}
