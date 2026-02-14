package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runEnvironment(dockerAppNames, terminalAppNames []string) error {
	fmt.Printf("Docker: %s | Terminal: %s\n\n",
		strings.Join(dockerAppNames, ", "), strings.Join(terminalAppNames, ", "))

	composeDir := filepath.Join(rootDir, "docker-compose")
	var composeFiles []string
	for _, f := range infraComposeFiles {
		composeFiles = append(composeFiles, filepath.Join(composeDir, f))
	}
	for _, name := range dockerAppNames {
		for _, app := range apps {
			if app.Name == name && app.ComposeFile != "" {
				composeFiles = append(composeFiles, filepath.Join(composeDir, app.ComposeFile))
			}
		}
	}

	// Clean up conflicting containers
	for _, c := range []string{"postgres", "postgres_instance", "redis", "redis_instance", "minio"} {
		exec.Command("docker", "rm", "-f", c).Run()
	}

	fmt.Println("==> Stopping existing stack...")
	_ = runDockerCompose(composeFiles, []string{"down", "--remove-orphans"}, true)

	for _, name := range terminalAppNames {
		for _, app := range apps {
			if app.Name == name {
				exec.Command("docker", "rm", "-f", app.ContainerName).Run()
			}
		}
	}

	fmt.Println("==> Starting Docker services...")
	if err := runDockerCompose(composeFiles, []string{"up", "-d"}, false); err != nil {
		return fmt.Errorf("failed to start Docker services: %w", err)
	}

	if len(terminalAppNames) > 0 {
		fmt.Println("==> Opening Terminal windows...")
		if err := openTerminalApps(terminalAppNames); err != nil {
			return fmt.Errorf("failed to open terminal apps: %w", err)
		}
	}
	fmt.Println("\n✅ All set!")
	return nil
}

func openTerminalApps(appNames []string) error {
	var osaArgs = []string{"-e", `tell application "Terminal"`}
	for _, name := range appNames {
		for _, app := range apps {
			if app.Name != name {
				continue
			}
			cmd := fmt.Sprintf("cd %q && cd %s && go run main.go", rootDir, app.BackendCmd)
			osaArgs = append(osaArgs, "-e", fmt.Sprintf(`do script "%s"`, escapeOsascript(cmd)))
			if app.FrontendCmd != "" {
				feCmd := fmt.Sprintf("cd %q && cd %s && npm run dev", rootDir, app.FrontendCmd)
				if app.Name == "dashboard" {
					feCmd = fmt.Sprintf("cd %q && cd %s && FILE_API_URL=http://localhost:8085 npm run dev", rootDir, app.FrontendCmd)
				}
				osaArgs = append(osaArgs, "-e", fmt.Sprintf(`do script "%s"`, escapeOsascript(feCmd)))
			}
		}
	}
	osaArgs = append(osaArgs, "-e", "activate", "-e", "end tell")
	cmd := exec.Command("osascript", osaArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
