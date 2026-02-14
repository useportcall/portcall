package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop all Portcall services",
		Long:  "Stop all Docker containers and clean up resources.",
		RunE:  stopCommand,
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be executed")
	return cmd
}

func stopCommand(cmd *cobra.Command, args []string) error {
	var err error
	rootDir, err = findRootDir()
	if err != nil {
		return fmt.Errorf("could not find project root: %w", err)
	}

	fmt.Println("===> Stopping all Portcall services...")
	composeDir := filepath.Join(rootDir, "docker-compose")

	// All compose files that might be running
	files := []string{
		filepath.Join(composeDir, "docker-compose.db.yml"),
		filepath.Join(composeDir, "docker-compose.auth.yml"),
		filepath.Join(composeDir, "docker-compose.apps.yml"),
		filepath.Join(composeDir, "docker-compose.tools.yml"),
		filepath.Join(composeDir, "docker-compose.workers.yml"),
		filepath.Join(composeDir, "docker-compose.dashboard.yml"),
		filepath.Join(composeDir, "docker-compose.checkout.yml"),
	}

	if err := runDockerCompose(files, []string{"down", "--remove-orphans"}, true); err != nil {
		return err
	}

	// Kill any remaining portcall containers
	containers := []string{
		"api", "admin", "dashboard", "checkout", "file-api",
		"quote", "postgres", "redis", "keycloak",
	}
	for _, c := range containers {
		exec.Command("docker", "rm", "-f", c).Run()
	}

	fmt.Println("✅ All services stopped")
	return nil
}
