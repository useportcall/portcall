package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var secretsOpts struct {
	cluster   string
	namespace string
}

func newSecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage Kubernetes secrets",
		Long:  "List, inspect, create, update, and delete Kubernetes secrets.",
	}
	cmd.PersistentFlags().StringVar(&secretsOpts.cluster, "cluster", "digitalocean", "Cluster name")
	cmd.PersistentFlags().StringVarP(&secretsOpts.namespace, "namespace", "n", "portcall", "Kubernetes namespace")
	cmd.AddCommand(newSecretsListCmd())
	cmd.AddCommand(newSecretsGetCmd())
	cmd.AddCommand(newSecretsSetCmd())
	cmd.AddCommand(newSecretsDeleteCmd())
	return cmd
}

func newSecretsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all secrets in the namespace",
		RunE:  runSecretsList,
	}
}

func newSecretsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <secret-name>",
		Short: "Show keys and values in a secret",
		Args:  cobra.ExactArgs(1),
		RunE:  runSecretsGet,
	}
}

func newSecretsSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <secret-name> <key>=<value> [key=value ...]",
		Short: "Create or update keys in a secret (preserves other keys)",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runSecretsSet,
	}
}

func newSecretsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <secret-name> [key ...]",
		Short: "Delete a secret or specific keys (with confirmation)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runSecretsDelete,
	}
}

func ensureSecretsCluster() error {
	return switchCluster(secretsOpts.cluster)
}

func runSecretsList(_ *cobra.Command, _ []string) error {
	if err := ensureSecretsCluster(); err != nil {
		return fmt.Errorf("switch cluster: %w", err)
	}
	return listSecrets(secretsOpts.namespace)
}

func runSecretsGet(_ *cobra.Command, args []string) error {
	if err := ensureSecretsCluster(); err != nil {
		return fmt.Errorf("switch cluster: %w", err)
	}
	return getSecret(secretsOpts.namespace, args[0])
}

func runSecretsSet(_ *cobra.Command, args []string) error {
	if err := ensureSecretsCluster(); err != nil {
		return fmt.Errorf("switch cluster: %w", err)
	}
	return setSecretKeys(secretsOpts.namespace, args[0], args[1:])
}

func runSecretsDelete(_ *cobra.Command, args []string) error {
	if err := ensureSecretsCluster(); err != nil {
		return fmt.Errorf("switch cluster: %w", err)
	}
	name := args[0]
	if len(args) == 1 {
		return deleteSecret(secretsOpts.namespace, name)
	}
	return deleteSecretKeys(secretsOpts.namespace, name, args[1:])
}
