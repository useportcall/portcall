package keycloak

import (
	"fmt"

	"github.com/spf13/cobra"
)

type switchClusterFn func(string) error

var opts struct {
	cluster   string
	namespace string
	realm     string
}

func NewCommand(switchCluster switchClusterFn) *cobra.Command {
	cmd := &cobra.Command{Use: "keycloak", Short: "Manage Keycloak configuration"}
	cmd.PersistentFlags().StringVar(&opts.cluster, "cluster", "digitalocean", "Cluster name")
	cmd.PersistentFlags().StringVarP(&opts.namespace, "namespace", "n", "portcall", "Kubernetes namespace")
	cmd.PersistentFlags().StringVar(&opts.realm, "realm", "dev", "Keycloak realm")
	cmd.AddCommand(&cobra.Command{
		Use:   "smtp-update",
		Short: "Update SMTP settings in the Keycloak realm",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSMTPUpdate(switchCluster)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "smtp-status",
		Short: "Show current SMTP and password-reset status",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSMTPStatus(switchCluster)
		},
	})
	return cmd
}

func runSMTPUpdate(switchCluster switchClusterFn) error {
	if err := switchCluster(opts.cluster); err != nil {
		return fmt.Errorf("switch cluster: %w", err)
	}
	kc, err := newClient(opts.namespace)
	if err != nil {
		return err
	}
	defer kc.close()
	return kc.updateSMTP(opts.realm)
}

func runSMTPStatus(switchCluster switchClusterFn) error {
	if err := switchCluster(opts.cluster); err != nil {
		return fmt.Errorf("switch cluster: %w", err)
	}
	kc, err := newClient(opts.namespace)
	if err != nil {
		return err
	}
	defer kc.close()
	return kc.showSMTPStatus(opts.realm)
}
