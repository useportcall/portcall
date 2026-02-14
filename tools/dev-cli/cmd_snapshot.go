package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSnapshotCmd() *cobra.Command {
	var (
		live   bool
		headed bool
		list   bool
		mode   string
	)

	cmd := &cobra.Command{
		Use:   "snapshot [targets...]",
		Short: "Take browser screenshots / videos and optionally send to Discord",
		Long: `Run Playwright snapshot tests for specific targets.

Targets can be individual snapshot names or group aliases.
By default nothing runs — you must specify at least one target,
a --mode filter, or --list to see what's available.

Examples:
  dev-cli snapshot --list                      # list available targets & groups
  dev-cli snapshot dashboard-home              # single target
  dev-cli snapshot dashboard-core checkout     # group + group
  dev-cli snapshot invoice-light --live        # upload to Discord
  dev-cli snapshot --mode fullscreen --live    # all fullscreen snapshots, live
  dev-cli snapshot dashboard-home --headed     # visible browser`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshot(args, list, live, headed, mode)
		},
	}

	cmd.Flags().BoolVar(&list, "list", false, "List available targets and groups")
	cmd.Flags().BoolVar(&live, "live", false, "Upload snapshots to Discord")
	cmd.Flags().BoolVar(&headed, "headed", false, "Run Playwright with a visible browser")
	cmd.Flags().StringVar(&mode, "mode", "", "Filter by mode: fullscreen, component, video")
	return cmd
}

func runSnapshot(args []string, list, live, headed bool, mode string) error {
	if list {
		printSnapshotList()
		return nil
	}
	targets, err := collectTargets(args, mode)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return fmt.Errorf("no targets specified; use --list to see available options")
	}
	return execSnapshot(targets, live, headed)
}
