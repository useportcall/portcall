package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/useportcall/portcall/tools/dev-cli/internal/deploy"
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra"
	"github.com/useportcall/portcall/tools/dev-cli/internal/keycloak"
	"github.com/useportcall/portcall/tools/dev-cli/internal/secrets"
	"github.com/useportcall/portcall/tools/dev-cli/internal/snapshot"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "portcall",
		Aliases: []string{"dev-cli"},
		Short:   "Development environment CLI for Portcall",
		Long:    "Manage Portcall infrastructure, deploys, and local development workflows.",
		RunE: func(*cobra.Command, []string) error {
			return runPortcallInteractive()
		},
	}

	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newE2ECmd())
	rootCmd.AddCommand(newSetupCmd())
	rootCmd.AddCommand(newStopCmd())
	rootCmd.AddCommand(deploy.NewCommand(findRootDir))
	rootCmd.AddCommand(secrets.NewCommand(switchCluster))
	rootCmd.AddCommand(keycloak.NewCommand(switchCluster))
	rootCmd.AddCommand(snapshot.NewCommand(findRootDir))
	rootCmd.AddCommand(infra.NewCommand(findRootDir))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
