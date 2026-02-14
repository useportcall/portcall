package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var skipDeps bool

func newSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "First-time setup for new developers",
		Long:  "Installs dependencies, builds SDK, starts Docker services, and verifies health.",
		RunE:  setupCommand,
	}
	cmd.Flags().BoolVar(&skipDeps, "skip-deps", false, "Skip npm/pnpm dependency installation")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be executed")
	return cmd
}

func setupCommand(cmd *cobra.Command, args []string) error {
	var err error
	rootDir, err = findRootDir()
	if err != nil {
		return fmt.Errorf("could not find project root: %w", err)
	}

	if err := checkPrerequisites(); err != nil {
		return err
	}

	if !skipDeps {
		if err := installDependencies(); err != nil {
			return err
		}
	}

	if err := startDockerServices(); err != nil {
		return err
	}

	if err := verifyServices(); err != nil {
		return err
	}

	printSetupComplete()
	return nil
}

func checkPrerequisites() error {
	fmt.Println("===> Checking prerequisites...")
	checks := []struct {
		cmd  string
		name string
		hint string
	}{
		{"docker", "Docker", "Install Docker Desktop"},
		{"pnpm", "pnpm", "Run: npm install -g pnpm"},
		{"go", "Go", "Install Go 1.20+"},
	}
	for _, c := range checks {
		if _, err := exec.LookPath(c.cmd); err != nil {
			return fmt.Errorf("❌ %s not found. %s", c.name, c.hint)
		}
	}
	if err := exec.Command("docker", "info").Run(); err != nil {
		return fmt.Errorf("❌ Docker daemon not running. Start Docker Desktop")
	}
	fmt.Println("✅ All prerequisites found")
	return nil
}

func installDependencies() error {
	fmt.Println("\n===> Installing dependencies...")
	if dryRun {
		fmt.Println("[DRY-RUN] pnpm install")
		fmt.Println("[DRY-RUN] cd sdks/node-typescript && pnpm build")
		return nil
	}
	c := exec.Command("pnpm", "install")
	c.Dir = rootDir
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("pnpm install failed: %w", err)
	}
	fmt.Println("\n===> Building SDK...")
	c = exec.Command("pnpm", "build")
	c.Dir = filepath.Join(rootDir, "sdks", "node-typescript")
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("SDK build failed: %w", err)
	}
	return nil
}
