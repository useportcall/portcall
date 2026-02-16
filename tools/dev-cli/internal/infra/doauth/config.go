package doauth

import (
	"os"
	"path/filepath"
	"strings"
)

func TokenFromDoctlConfig(context string) (string, error) {
	cfg := strings.TrimSpace(os.Getenv("DOCTL_CONFIG"))
	if cfg == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cfg = filepath.Join(home, "Library", "Application Support", "doctl", "config.yaml")
	}
	buf, err := os.ReadFile(cfg)
	if err != nil {
		return "", err
	}
	return FindContextToken(string(buf), context), nil
}

func FindContextToken(cfg string, context string) string {
	in := false
	for _, line := range strings.Split(cfg, "\n") {
		if strings.TrimSpace(line) == "auth-contexts:" {
			in = true
			continue
		}
		if !in {
			continue
		}
		if len(line) > 0 && line[0] != ' ' {
			break
		}
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) != context {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[1]), "\"")
	}
	return ""
}
