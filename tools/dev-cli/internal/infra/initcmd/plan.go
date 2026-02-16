package initcmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

func buildPlan(opts Options, deps Deps) (Plan, error) {
	providers, err := deps.ResolveProviders(opts.Providers, opts.Provider)
	if err != nil {
		return Plan{}, err
	}
	if opts.Mode != "micro" {
		return Plan{}, fmt.Errorf("unsupported mode %q (supported: micro)", opts.Mode)
	}
	if opts.Action != "create" && opts.Action != "update" {
		return Plan{}, fmt.Errorf("unsupported action %q (supported: create|update)", opts.Action)
	}
	if opts.Step != "all" && opts.Step != "core" && opts.Step != "services" {
		return Plan{}, fmt.Errorf("unsupported step %q (supported: all|core|services)", opts.Step)
	}
	if strings.TrimSpace(opts.Cluster) == "" {
		return Plan{}, fmt.Errorf("--cluster alias is required")
	}
	if opts.NodeCount < 1 {
		return Plan{}, fmt.Errorf("--node-count must be >= 1")
	}
	if strings.TrimSpace(opts.RedisSize) == "" {
		return Plan{}, fmt.Errorf("--redis-size is required")
	}
	ips, err := deps.ResolveAllowedIPs(opts.AllowedIPs, opts.AllowedIPsFrom)
	if err != nil {
		return Plan{}, err
	}
	registry := fmt.Sprintf("registry.digitalocean.com/%s", strings.TrimSpace(opts.RegistryName))
	return Plan{
		Alias:        strings.TrimSpace(opts.Cluster),
		Providers:    providers,
		Provider:     providers[0],
		Action:       strings.TrimSpace(opts.Action),
		Step:         strings.TrimSpace(opts.Step),
		ClusterName:  strings.TrimSpace(opts.ClusterName),
		Region:       strings.TrimSpace(opts.Region),
		NodeSize:     strings.TrimSpace(opts.NodeSize),
		NodeCount:    opts.NodeCount,
		RedisSize:    strings.TrimSpace(opts.RedisSize),
		RegistryName: strings.TrimSpace(opts.RegistryName),
		Registry:     registry,
		VPCCIDR:      strings.TrimSpace(opts.VPCCIDR),
		SpacesRegion: strings.TrimSpace(opts.SpacesRegion),
		SpacesPrefix: strings.TrimSpace(opts.SpacesPrefix),
		Domain:       strings.TrimSpace(opts.Domain),
		AllowedIPs:   ips,
		Mode:         opts.Mode,
		Namespace:    "portcall",
		StackDir:     filepath.Join(deps.RootDir(), "infra", "digitalocean", "terraform"),
	}, nil
}
