package cleanupcmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Deps struct {
	EnsureRootDir func() error
	SwitchCluster func(cluster string) error
	ReadInput     func(prompt string) string
	RunCmd        func(name string, args ...string) error
	Warn          func(msg string, args ...any)
	Plain         func(msg string, args ...any)
	OK            func(msg string, args ...any)
}

func New(deps Deps) *cobra.Command {
	var opts struct {
		cluster   string
		namespace string
		yes       bool
		dryRun    bool
	}
	cmd := &cobra.Command{Use: "cleanup", Short: "Cleanup legacy infrastructure artifacts"}
	legacy := &cobra.Command{
		Use:   "legacy",
		Short: "Remove deprecated in-cluster resources (smtp-relay, old jobs)",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := deps.EnsureRootDir(); err != nil {
				return err
			}
			if err := deps.SwitchCluster(opts.cluster); err != nil {
				return err
			}
			resources := [][]string{{"deployment", "smtp-relay"}, {"service", "smtp-relay"}, {"configmap", "smtp-relay-config"}, {"secret", "smtp-relay-secrets"}}
			if !opts.yes && !opts.dryRun {
				deps.Warn("Will cleanup legacy resources in namespace %s", opts.namespace)
				for _, resource := range resources {
					deps.Plain("- %s/%s", resource[0], resource[1])
				}
				if strings.ToLower(deps.ReadInput("Proceed? [y/N]: ")) != "y" {
					return fmt.Errorf("cleanup canceled")
				}
			}
			for _, resource := range resources {
				if opts.dryRun {
					deps.Plain("[DRY-RUN] kubectl delete %s %s -n %s --ignore-not-found", resource[0], resource[1], opts.namespace)
					continue
				}
				if err := deps.RunCmd("kubectl", "delete", resource[0], resource[1], "-n", opts.namespace, "--ignore-not-found"); err != nil {
					return fmt.Errorf("cleanup %s/%s: %w", resource[0], resource[1], err)
				}
			}
			deps.OK("Legacy cleanup complete")
			return nil
		},
	}
	legacy.Flags().StringVar(&opts.cluster, "cluster", "digitalocean", "Cluster alias")
	legacy.Flags().StringVarP(&opts.namespace, "namespace", "n", "portcall", "Kubernetes namespace")
	legacy.Flags().BoolVarP(&opts.yes, "yes", "y", false, "Skip confirmation prompt")
	legacy.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Print actions without deleting")
	cmd.AddCommand(legacy)
	return cmd
}
