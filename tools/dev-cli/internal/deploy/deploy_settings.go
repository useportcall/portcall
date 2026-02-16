package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type deploySettings struct {
	SkipMigration bool `json:"skipMigration"`
}

func deploySettingsPath(root string) string {
	return filepath.Join(root, ".dev-cli.deploy.json")
}

func loadDeploySettings(root string) (deploySettings, error) {
	var settings deploySettings
	data, err := os.ReadFile(deploySettingsPath(root))
	if errors.Is(err, os.ErrNotExist) {
		return settings, nil
	}
	if err != nil {
		return settings, fmt.Errorf("failed to read deploy settings: %w", err)
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		return settings, fmt.Errorf("failed to parse deploy settings: %w", err)
	}
	return settings, nil
}

func saveDeploySettings(root string, settings deploySettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode deploy settings: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(deploySettingsPath(root), data, 0o644); err != nil {
		return fmt.Errorf("failed to write deploy settings: %w", err)
	}
	return nil
}
