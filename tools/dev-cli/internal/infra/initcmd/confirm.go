package initcmd

import "fmt"

func confirmInit(plan Plan, opts Options, deps Deps) error {
	if opts.Yes || opts.DryRun {
		return nil
	}
	deps.Section("Infra Init Summary")
	deps.Plain("Alias       : %s", plan.Alias)
	deps.Plain("Providers   : %v", plan.Providers)
	deps.Plain("Action      : %s", plan.Action)
	deps.Plain("Cluster     : %s", plan.ClusterName)
	deps.Plain("Node pool   : %s x%d", plan.NodeSize, plan.NodeCount)
	deps.Plain("Redis size  : %s", plan.RedisSize)
	deps.Plain("Step        : %s", plan.Step)
	deps.Plain("Region      : %s", plan.Region)
	deps.Plain("Mode        : %s", plan.Mode)
	deps.Plain("Registry    : %s", plan.Registry)
	deps.Plain("Domain      : %s", plan.Domain)
	deps.Plain("DNS provider: %s", opts.DNSProvider)
	deps.Plain("DNS auto    : %t", opts.DNSAuto)
	deps.Plain("Allowed IPs : %v", plan.AllowedIPs)
	deps.Plain("Dry run     : %t", opts.DryRun)
	if deps.ReadInput("Proceed with infra init? [y/N]: ") != "y" {
		return fmt.Errorf("infra init canceled")
	}
	return nil
}

func previewDryRun(plan Plan, opts Options, steps []Step, deps Deps) {
	deps.Section("Dry-Run Preview")
	deps.Plain("Alias      : %s", plan.Alias)
	deps.Plain("Provider   : %s", plan.Provider)
	deps.Plain("Action     : %s", plan.Action)
	deps.Plain("Mode/Step  : %s / %s", plan.Mode, plan.Step)
	deps.Plain("Cluster    : %s (%s, %s x%d)", plan.ClusterName, plan.Region, plan.NodeSize, plan.NodeCount)
	deps.Plain("Redis size : %s", plan.RedisSize)
	deps.Plain("Registry   : %s", plan.Registry)
	deps.Plain("Domain     : %s", plan.Domain)
	deps.Plain("DNS        : provider=%s auto=%t", opts.DNSProvider, opts.DNSAuto)
	deps.Plain("Namespace  : %s", plan.Namespace)
	if len(steps) == 0 {
		deps.Warn("No terraform steps selected in this preview")
	} else {
		deps.Info("Terraform steps that would run:")
		for idx, step := range steps {
			deps.Plain("  %d. %s (%d target(s))", idx+1, step.Name, len(step.Targets))
		}
	}
	deps.Plain("No terraform, doctl, kubectl, or helm command was executed in dry-run mode.")
	deps.Plain("Try next dry-runs:")
	deps.Plain("- go run ./tools/dev-cli infra doctor --cluster %s", plan.Alias)
	deps.Plain("- go run ./tools/dev-cli infra status --cluster %s", plan.Alias)
	deps.Plain("- go run ./tools/dev-cli infra cleanup legacy --cluster %s --dry-run", plan.Alias)
}
