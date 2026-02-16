package secrets

import (
	"fmt"

	"github.com/spf13/cobra"
)

type switchClusterFn func(string) error

var opts struct {
	cluster   string
	namespace string
}

func NewCommand(switchCluster switchClusterFn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage Kubernetes secrets",
		Long:  "List, inspect, create, update, and delete Kubernetes secrets.",
	}
	cmd.PersistentFlags().StringVar(&opts.cluster, "cluster", "digitalocean", "Cluster name")
	cmd.PersistentFlags().StringVarP(&opts.namespace, "namespace", "n", "portcall", "Kubernetes namespace")
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List all secrets in the namespace", RunE: runList(switchCluster)})
	cmd.AddCommand(&cobra.Command{Use: "get <secret-name>", Short: "Show keys and values in a secret", Args: cobra.ExactArgs(1), RunE: runGet(switchCluster)})
	cmd.AddCommand(&cobra.Command{
		Use:   "set <secret-name> <key>=<value> [key=value ...]",
		Short: "Create or update keys in a secret (preserves other keys)",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runSet(switchCluster),
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "delete <secret-name> [key ...]",
		Short: "Delete a secret or specific keys (with confirmation)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runDelete(switchCluster),
	})
	return cmd
}

func runList(s switchClusterFn) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		if err := s(opts.cluster); err != nil {
			return fmt.Errorf("switch cluster: %w", err)
		}
		return list(opts.namespace)
	}
}

func runGet(s switchClusterFn) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if err := s(opts.cluster); err != nil {
			return fmt.Errorf("switch cluster: %w", err)
		}
		return get(opts.namespace, args[0])
	}
}

func runSet(s switchClusterFn) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if err := s(opts.cluster); err != nil {
			return fmt.Errorf("switch cluster: %w", err)
		}
		return setKeys(opts.namespace, args[0], args[1:])
	}
}

func runDelete(s switchClusterFn) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if err := s(opts.cluster); err != nil {
			return fmt.Errorf("switch cluster: %w", err)
		}
		if len(args) == 1 {
			return deleteSecret(opts.namespace, args[0])
		}
		return deleteKeys(opts.namespace, args[0], args[1:])
	}
}
