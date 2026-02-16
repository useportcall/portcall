package deploy

import "github.com/spf13/cobra"

type deployOptions struct {
	appInput       string
	cluster        string
	version        string
	valuesFile     string
	skipBuild      bool
	skipMigration  bool
	skipSmokeTests bool
	skipPreflight  bool
	runUnitTests   bool
	runIntegTests  bool
	runE2ETests    bool
	listApps       bool
	dryRun         bool
	yes            bool
}

var deployOpts deployOptions

func newDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy selected apps to Kubernetes",
		RunE:  runDeploy,
	}
	cmd.Flags().StringVar(&deployOpts.appInput, "apps", "", "Comma-separated app names/numbers or 'all'")
	cmd.Flags().StringVar(&deployOpts.cluster, "cluster", "digitalocean", "Cluster name")
	cmd.Flags().StringVar(&deployOpts.version, "version", "", "Version bump: major|minor|patch|skip")
	cmd.Flags().StringVar(&deployOpts.valuesFile, "values", "", "Values file (default: k8s/deploy/<cluster>/values.yaml)")
	cmd.Flags().BoolVar(&deployOpts.skipBuild, "skip-build", false, "Skip Docker build and push")
	cmd.Flags().BoolVar(&deployOpts.skipMigration, "skip-migration", false, "Skip database migration (saved as default when explicitly set)")
	cmd.Flags().BoolVar(&deployOpts.skipSmokeTests, "skip-tests", false, "Skip smoke tests")
	cmd.Flags().BoolVar(&deployOpts.skipPreflight, "skip-preflight", false, "Skip pre-deploy tests")
	cmd.Flags().BoolVar(&deployOpts.runUnitTests, "unit-tests", true, "Run unit tests before deploy")
	cmd.Flags().BoolVar(&deployOpts.runIntegTests, "integration-tests", true, "Run integration tests before deploy")
	cmd.Flags().BoolVar(&deployOpts.runE2ETests, "e2e-tests", false, "Run e2e tests before deploy")
	cmd.Flags().BoolVar(&deployOpts.listApps, "list-apps", false, "List available apps and exit")
	cmd.Flags().BoolVar(&deployOpts.dryRun, "dry-run", false, "Print actions without making changes")
	cmd.Flags().BoolVarP(&deployOpts.yes, "yes", "y", false, "Skip confirmation prompts")
	return cmd
}
