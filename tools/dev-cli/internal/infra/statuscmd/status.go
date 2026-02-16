package statuscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	infraState "github.com/useportcall/portcall/tools/dev-cli/internal/infra/state"
)

type Deps struct {
	EnsureRootDir         func() error
	GetInfraClusterState  func(cluster string) (infraState.ClusterState, bool)
	ResolveClusterContext func(cluster string) string
	Warn                  func(msg string, args ...any)
	Plain                 func(msg string, args ...any)
}

func New(deps Deps) *cobra.Command {
	var cluster string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show saved infra state for a cluster alias",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := deps.EnsureRootDir(); err != nil {
				return err
			}
			cfg, ok := deps.GetInfraClusterState(cluster)
			if !ok {
				deps.Warn("No saved infra state for alias %s", cluster)
				deps.Plain("Run: go run ./tools/dev-cli infra pull --cluster %s", cluster)
				return nil
			}
			deps.Plain("Alias:     %s", cluster)
			deps.Plain("Provider:  %s", cfg.Provider)
			deps.Plain("Action:    %s", cfg.InitAction)
			deps.Plain("Context:   %s", cfg.Context)
			deps.Plain("Cluster:   %s", cfg.Cluster)
			deps.Plain("Namespace: %s", cfg.Namespace)
			deps.Plain("Mode:      %s", cfg.Mode)
			deps.Plain("Registry:  %s", cfg.Registry)
			deps.Plain("Values:    %s", cfg.Values)
			deps.Plain("Context default resolver: %s", deps.ResolveClusterContext(cluster))
			deps.Plain("Audit full cluster wiring: go run ./tools/dev-cli infra doctor --cluster %s", cluster)
			if cfg.Context == "" && cfg.Values == "" {
				return fmt.Errorf("saved state is incomplete for alias %s", cluster)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&cluster, "cluster", "digitalocean", "Cluster alias")
	return cmd
}
