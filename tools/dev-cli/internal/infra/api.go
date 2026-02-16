package infra

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand(rootFind func() (string, error)) *cobra.Command {
	SetRootResolver(rootFind)
	return newInfraCmd()
}

func ResolveClusterContext(root, cluster string) string {
	cfg, ok := LoadClusterState(root, cluster)
	if ok && strings.TrimSpace(cfg.Context) != "" {
		return cfg.Context
	}
	if cluster == "digitalocean" {
		return "do-k8s-portcall-prod"
	}
	return cluster
}

func ResolveValuesFile(root, cluster, current string) string {
	if strings.TrimSpace(current) != "" {
		return current
	}
	cfg, ok := LoadClusterState(root, cluster)
	if ok && strings.TrimSpace(cfg.Values) != "" {
		return cfg.Values
	}
	return fmt.Sprintf("%s/k8s/deploy/%s/values.yaml", root, cluster)
}

func ResolveImageRepository(root, cluster, defaultRepo string) string {
	cfg, ok := LoadClusterState(root, cluster)
	if !ok || strings.TrimSpace(cfg.Registry) == "" {
		return defaultRepo
	}
	return strings.TrimSuffix(cfg.Registry, "/") + "/" + path.Base(defaultRepo)
}
