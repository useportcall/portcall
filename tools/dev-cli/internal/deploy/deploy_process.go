package deploy

import (
	"fmt"
	"strings"
	"time"
)

func processOneApp(name string, versions map[string]string) error {
	app, okA := findDeployApp(name)
	if !okA {
		return fmt.Errorf("app not found: %s", name)
	}
	warn("====== Processing: %s ======", name)
	ver := "skip"
	if !app.ChartOnly {
		cur := currentVersion(name)
		ver = cur
		if deployOpts.version != "skip" {
			v, err := bumpVersion(cur, deployOpts.version)
			if err != nil {
				return err
			}
			ver = v
		}
	}
	ok("Using version: %s", ver)
	if !deployOpts.skipBuild {
		if err := buildAndPush(app, ver, deployOpts.dryRun); err != nil {
			return err
		}
	}
	if !app.ChartOnly {
		versions[helmValueName(name)] = ver
	}
	if deployOpts.dryRun {
		plain("[DRY-RUN] Would deploy %s via Helm", name)
		return nil
	}
	clearStuckRelease()
	return helmDeploy(name, ver, versions)
}

func helmDeploy(app, version string, versions map[string]string) error {
	target, _ := findDeployApp(app)
	args := []string{"upgrade", "--install", "portcall", "./k8s/portcall-chart", "-f", deployOpts.valuesFile, "--namespace", "portcall", "--create-namespace"}
	args = append(args, versionSetArgs(versions)...)
	args = append(args, "--set", fmt.Sprintf("%s.enabled=true", helmValueName(app)), "--timeout=2m")
	if err := runCmd("helm", args...); err != nil {
		return fmt.Errorf("helm upgrade failed for %s: %w", app, err)
	}
	if target.ChartOnly && app == "observability" {
		return waitForObservability()
	}
	info("Watching rollout for %s...", deployName(app))
	if watchRollout(deployName(app), 60*time.Second) {
		ok("Successfully deployed %s:%s", app, version)
		return nil
	}
	out, _ := runCmdOut("kubectl", "get", "deployment", deployName(app), "-n", "portcall", "-o", "jsonpath={.status.readyReplicas}/{.spec.replicas}")
	warn("Rollout timeout, status: %s", strings.TrimSpace(out))
	if strings.ToLower(readInput("Continue anyway? [y/N]: ")) != "y" {
		return fmt.Errorf("deployment aborted by user")
	}
	return nil
}
