package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"
)

func startDockerServices() error {
	fmt.Println("\n===> Starting Docker services...")
	composeDir := filepath.Join(rootDir, "docker-compose")
	files := []string{
		filepath.Join(composeDir, "docker-compose.db.yml"),
		filepath.Join(composeDir, "docker-compose.auth.yml"),
		filepath.Join(composeDir, "docker-compose.apps.yml"),
	}
	// Clean up conflicting containers
	for _, c := range []string{"postgres", "postgres_instance", "redis", "redis_instance"} {
		exec.Command("docker", "rm", "-f", c).Run()
	}
	if err := runDockerCompose(files, []string{"up", "-d"}, false); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}
	fmt.Println("Waiting for services to initialize (30s)...")
	if !dryRun {
		time.Sleep(30 * time.Second)
	}
	return nil
}

func verifyServices() error {
	fmt.Println("\n===> Verifying services...")
	if dryRun {
		fmt.Println("[DRY-RUN] curl http://localhost:9080/ping")
		return nil
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:9080/ping")
	if err != nil {
		return fmt.Errorf("❌ API not responding: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("❌ API returned status %d", resp.StatusCode)
	}
	fmt.Println("✅ API is healthy")
	return nil
}

func printSetupComplete() {
	fmt.Println("\n=== ✨ Setup Complete! ===")
	fmt.Println("\nServices running:")
	fmt.Println("  • Dashboard: http://localhost:8082")
	fmt.Println("  • API:       http://localhost:9080")
	fmt.Println("  • Admin:     http://localhost:8081")
	fmt.Println("\nNext: ./dev-cli run --preset=dashboard")
}
