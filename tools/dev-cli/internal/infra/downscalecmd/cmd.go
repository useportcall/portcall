package downscalecmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func New(deps Deps) *cobra.Command {
	opts := Options{Cluster: "digitalocean", ClusterName: "portcall-prod", NodeCount: 1, RedisSize: "db-s-1vcpu-1gb", SmokeTimeout: 300}
	cmd := &cobra.Command{Use: "downscale", Short: "Interactively downscale cluster + Redis for playground usage", RunE: func(c *cobra.Command, _ []string) error {
		plan, err := BuildPlan(c, opts, deps)
		if err != nil {
			return err
		}
		PrintSummary(plan, deps)
		if !opts.Yes && !deps.AskYesNo("Proceed with downscale? [y/N]: ", false) {
			return fmt.Errorf("downscale canceled")
		}
		if opts.DryRun {
			deps.OK("Downscale dry-run complete (no changes applied)")
			return nil
		}
		if err := RunApply(plan, deps); err != nil {
			return err
		}
		if !opts.SkipSmokeCheck {
			if err := RunSmokeCheck(opts.Cluster, opts.SmokeTimeout, deps); err != nil {
				return err
			}
		}
		deps.OK("Downscale complete")
		return nil
	}}
	cmd.Flags().StringVar(&opts.Cluster, "cluster", "digitalocean", "Cluster alias")
	cmd.Flags().StringVar(&opts.ClusterName, "name", "portcall-prod", "Kubernetes cluster name")
	cmd.Flags().StringVar(&opts.NodeSize, "node-size", "", "Kubernetes node size slug")
	cmd.Flags().IntVar(&opts.NodeCount, "node-count", 1, "Kubernetes node count")
	cmd.Flags().StringVar(&opts.RedisSize, "redis-size", "db-s-1vcpu-1gb", "Managed Redis size slug")
	cmd.Flags().StringVar(&opts.PostgresSize, "postgres-size", "", "Managed Postgres size slug (optional)")
	cmd.Flags().BoolVar(&opts.SkipSmokeCheck, "skip-smoke-check", false, "Skip post-change smoke checks")
	cmd.Flags().IntVar(&opts.SmokeTimeout, "smoke-timeout-sec", 300, "Smoke check timeout in seconds")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Show actions without applying changes")
	return cmd
}
