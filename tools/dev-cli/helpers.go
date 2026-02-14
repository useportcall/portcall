package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	parallel int
	dryRun   bool
	rootDir  string
)

func findRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	exe, err := os.Executable()
	if err == nil {
		dir = filepath.Dir(exe)
		for i := 0; i < 5; i++ {
			if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
				return dir, nil
			}
			dir = filepath.Dir(dir)
		}
	}
	return "", fmt.Errorf("could not find go.work in any parent directory")
}

func parseAppList(input string) []string {
	var result []string
	for _, s := range strings.Split(input, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func validateApps(appNames []string) error {
	for _, name := range appNames {
		found := false
		for _, app := range apps {
			if app.Name == name {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown app: %s (use --list to see available apps)", name)
		}
	}
	return nil
}

func runDockerCompose(composeFiles, args []string, allowFail bool) error {
	var cmdArgs []string
	for _, f := range composeFiles {
		cmdArgs = append(cmdArgs, "-f", f)
	}
	cmdArgs = append(cmdArgs, args...)
	if dryRun {
		fmt.Printf("[DRY-RUN] docker compose %s\n", strings.Join(cmdArgs, " "))
		return nil
	}
	cmd := exec.Command("docker", append([]string{"compose"}, cmdArgs...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("COMPOSE_PARALLEL_LIMIT=%d", parallel))
	err := cmd.Run()
	if err != nil && !allowFail {
		return err
	}
	return nil
}

func escapeOsascript(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
