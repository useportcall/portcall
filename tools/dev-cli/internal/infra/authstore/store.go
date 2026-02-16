package authstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type State struct {
	CloudflareAPIToken string `json:"cloudflareApiToken,omitempty"`
}

func path(root string) string {
	return filepath.Join(root, ".dev-cli.auth.json")
}

func Load(root string) (State, error) {
	var out State
	buf, err := os.ReadFile(path(root))
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("read auth store: %w", err)
	}
	if err := json.Unmarshal(buf, &out); err != nil {
		return out, fmt.Errorf("parse auth store: %w", err)
	}
	return out, nil
}

func Save(root string, state State) error {
	buf, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode auth store: %w", err)
	}
	buf = append(buf, '\n')
	if err := os.WriteFile(path(root), buf, 0o600); err != nil {
		return fmt.Errorf("write auth store: %w", err)
	}
	return nil
}

func SaveCloudflareToken(root string, token string) error {
	state, err := Load(root)
	if err != nil {
		return err
	}
	state.CloudflareAPIToken = strings.TrimSpace(token)
	return Save(root, state)
}
