package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dev-cli",
		Short: "Development environment CLI for Portcall",
		Long:  "Manage the Portcall dev environment and run e2e tests.",
	}

	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newE2ECmd())
	rootCmd.AddCommand(newSetupCmd())
	rootCmd.AddCommand(newStopCmd())
	rootCmd.AddCommand(newDeployCmd())
	rootCmd.AddCommand(newSecretsCmd())
	rootCmd.AddCommand(newSnapshotCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
