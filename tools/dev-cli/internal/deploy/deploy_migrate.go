package deploy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// runLocalMigration fetches DATABASE_URL from the K8s secret,
// then runs `go run ./migrate` locally against the remote DB.
func runLocalMigration() error {
	info("Running database migration locally...")
	dbURL, err := fetchDatabaseURL()
	if err != nil {
		return fmt.Errorf("failed to get DATABASE_URL: %w", err)
	}
	if deployOpts.dryRun {
		plain("[DRY-RUN] Would run migration against remote DB")
		return nil
	}
	migrateDir := filepath.Join(rootDir, "migrate")
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = migrateDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"DATABASE_URL="+dbURL,
		"APP_ENV=production",
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	ok("Database migration completed successfully")
	return nil
}

// fetchDatabaseURL reads DATABASE_URL from the portcall-secrets K8s secret.
func fetchDatabaseURL() (string, error) {
	out, err := runCmdOut("kubectl", "get", "secret", "portcall-secrets",
		"-n", "portcall", "-o", "json")
	if err != nil {
		return "", fmt.Errorf("kubectl get secret failed: %w", err)
	}
	var secret struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &secret); err != nil {
		return "", fmt.Errorf("failed to parse secret JSON: %w", err)
	}
	encoded, ok := secret.Data["DATABASE_URL"]
	if !ok {
		return "", fmt.Errorf("DATABASE_URL not found in secret")
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode DATABASE_URL: %w", err)
	}
	return string(decoded), nil
}
