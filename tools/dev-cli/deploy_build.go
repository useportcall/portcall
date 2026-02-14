package main

import (
	"fmt"
	"strings"
)

func buildAndPush(app deployApp, version string, dry bool) error {
	if app.ChartOnly {
		info("Chart-only target %s: skipping build/push", app.Name)
		return nil
	}
	if dry {
		warn("[DRY-RUN] Would build and push %s:%s", app.Image, version)
		return nil
	}
	info("Building Docker image for %s...", app.Name)
	args := []string{"build", "--platform", "linux/amd64", "-t", app.Image + ":" + version, "-t", app.Image + ":latest", "-f", app.Dockerfile}
	if app.Name == "admin" {
		args = append(args, "--build-arg", "VITE_KEYCLOAK_URL=https://auth.useportcall.com")
	}
	args = append(args, app.Context)
	if hasBuildx() {
		args = append([]string{"buildx", "build"}, args[1:]...)
		args = append(args, "--load")
	}
	if err := runCmd("docker", args...); err != nil {
		return fmt.Errorf("build failed for %s: %w", app.Name, err)
	}
	info("Pushing image to registry...")
	if err := runCmd("docker", "push", app.Image+":"+version); err != nil {
		return fmt.Errorf("push failed for %s: %w", app.Name, err)
	}
	_ = runCmd("docker", "push", app.Image+":latest")
	ok("Image pushed: %s:%s", app.Image, version)
	return nil
}

func resolveValuesFile(cluster, current string) string {
	if strings.TrimSpace(current) != "" {
		return current
	}
	return fmt.Sprintf("%s/k8s/deploy/%s/values.yaml", rootDir, cluster)
}
