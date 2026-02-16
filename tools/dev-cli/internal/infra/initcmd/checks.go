package initcmd

import (
	"fmt"
	"strings"
)

func confirmGuardrails(plan Plan, opts Options, deps Deps) error {
	signals := existingClusterSignals(plan, deps)
	if len(signals) == 0 {
		return nil
	}
	deps.Section("Existing Cluster Detected")
	deps.Warn("Found existing cluster signals for %s", plan.ClusterName)
	for _, s := range signals {
		deps.Plain("- %s", s)
	}
	if opts.Action != "create" {
		if opts.Yes {
			return nil
		}
		deps.Info("Default for existing cluster with update action: continue")
		if !deps.AskYesNo("Proceed with update action? [y/N]: ", false) {
			return fmt.Errorf("init canceled")
		}
		return nil
	}
	deps.Warn("Action is create; continuing may overwrite existing resources")
	if opts.Yes {
		return fmt.Errorf("refusing create action with --yes on existing cluster; use --action update or rerun interactively")
	}
	token := "overwrite " + plan.ClusterName
	deps.Info("Safety confirmation required to continue create mode")
	if deps.AskText("Type \""+token+"\" to continue: ", "") != token {
		return fmt.Errorf("init canceled")
	}
	if !deps.AskYesNo("Final confirmation: proceed with create on existing cluster? [y/N]: ", false) {
		return fmt.Errorf("init canceled")
	}
	return nil
}

func existingClusterSignals(plan Plan, deps Deps) []string {
	signals := []string{}
	if cfg, ok := deps.GetClusterState(plan.Alias); ok {
		signals = append(signals, fmt.Sprintf("alias %s exists in .dev-cli.infra.json", plan.Alias))
		if cfg.Cluster != "" {
			signals = append(signals, fmt.Sprintf("saved cluster name: %s", cfg.Cluster))
		}
	}
	out, err := deps.RunCmdOut("doctl", "kubernetes", "cluster", "get", plan.ClusterName, "--format", "Name", "--no-header")
	if err == nil && strings.TrimSpace(out) != "" {
		signals = append(signals, fmt.Sprintf("digitalocean cluster exists: %s", strings.TrimSpace(out)))
	}
	return signals
}
