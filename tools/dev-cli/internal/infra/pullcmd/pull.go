package pullcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	infraState "github.com/useportcall/portcall/tools/dev-cli/internal/infra/state"
)

type Deps struct {
	EnsureRootDir        func() error
	SwitchCluster        func(cluster string) error
	RunCmdOut            func(name string, args ...string) (string, error)
	SetInfraClusterState func(cluster string, cfg infraState.ClusterState) error
	GetInfraClusterState func(cluster string) (infraState.ClusterState, bool)
	RootDir              func() string
	OK                   func(msg string, args ...any)
	Plain                func(msg string, args ...any)
}

func New(deps Deps) *cobra.Command {
	var opts struct {
		cluster      string
		namespace    string
		valuesSource string
		clusterName  string
		mode         string
	}
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull running cluster state into local dev-cli infra state",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := deps.EnsureRootDir(); err != nil {
				return err
			}
			if err := deps.SwitchCluster(opts.cluster); err != nil {
				return err
			}
			ctx, err := deps.RunCmdOut("kubectl", "config", "current-context")
			if err != nil {
				return err
			}
			current, _ := deps.GetInfraClusterState(opts.cluster)
			registry := firstNonEmpty(pullRegistry(deps.RunCmdOut, opts.namespace), current.Registry)
			domain := pullDomain(deps.RunCmdOut, opts.namespace)
			clusterName := firstNonEmpty(opts.clusterName, firstNonEmpty(pullClusterName(deps.RunCmdOut), current.Cluster))
			valuesPath, err := copyPulledValues(deps.RootDir(), opts.valuesSource, opts.cluster)
			if err != nil {
				return err
			}
			if err := deps.SetInfraClusterState(opts.cluster, infraState.ClusterState{
				Context: ctx, Registry: registry, Values: valuesPath, Mode: opts.mode,
				Cluster: clusterName, Namespace: opts.namespace, Provider: firstNonEmpty(current.Provider, "digitalocean"),
			}); err != nil {
				return err
			}
			deps.OK("Pulled cluster state for alias %s", opts.cluster)
			deps.Plain("Context:  %s", ctx)
			deps.Plain("Registry: %s", registry)
			deps.Plain("Domain:   %s", domain)
			deps.Plain("Values:   %s", valuesPath)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.cluster, "cluster", "digitalocean", "Cluster alias")
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "portcall", "Kubernetes namespace")
	cmd.Flags().StringVar(&opts.valuesSource, "values-source", "k8s/deploy/digitalocean/values.yaml", "Source values file to align with")
	cmd.Flags().StringVar(&opts.clusterName, "name", "", "Cluster name override for saved state")
	cmd.Flags().StringVar(&opts.mode, "mode", "micro", "Cluster mode label to store")
	return cmd
}

func copyPulledValues(root, src, alias string) (string, error) {
	path := filepath.Join(root, src)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read values source: %w", err)
	}
	dir := filepath.Join(root, ".infra", alias)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create pulled values dir: %w", err)
	}
	dst := filepath.Join(dir, "values.pulled.yaml")
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return "", fmt.Errorf("write pulled values: %w", err)
	}
	return dst, nil
}

func pullRegistry(runCmdOut func(name string, args ...string) (string, error), namespace string) string {
	img, err := runCmdOut("kubectl", "get", "deployment", "api", "-n", namespace, "-o", "jsonpath={.spec.template.spec.containers[0].image}")
	if err != nil || img == "" {
		return ""
	}
	repo := strings.Split(img, ":")[0]
	parts := strings.Split(repo, "/")
	if len(parts) < 2 {
		return ""
	}
	return strings.Join(parts[:len(parts)-1], "/")
}

func pullDomain(runCmdOut func(name string, args ...string) (string, error), namespace string) string {
	host, err := runCmdOut("kubectl", "get", "ingress", "portcall-ingress", "-n", namespace, "-o", "jsonpath={.spec.rules[0].host}")
	if err != nil || host == "" {
		return ""
	}
	parts := strings.Split(host, ".")
	if len(parts) < 3 {
		return host
	}
	return strings.Join(parts[1:], ".")
}

func pullClusterName(runCmdOut func(name string, args ...string) (string, error)) string {
	name, err := runCmdOut("kubectl", "get", "nodes", "-o", "jsonpath={.items[0].metadata.labels.doks\\.digitalocean\\.com/cluster-name}")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(name)
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
