package snapshot

import (
	"fmt"

	"github.com/spf13/cobra"
)

type rootFinder func() (string, error)

func NewCommand(findRoot rootFinder) *cobra.Command {
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
By default nothing runs. Specify at least one target,
a --mode filter, or --list to see what's available.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(findRoot, args, list, live, headed, mode)
		},
	}
	cmd.Flags().BoolVar(&list, "list", false, "List available targets and groups")
	cmd.Flags().BoolVar(&live, "live", false, "Upload snapshots to Discord")
	cmd.Flags().BoolVar(&headed, "headed", false, "Run Playwright with a visible browser")
	cmd.Flags().StringVar(&mode, "mode", "", "Filter by mode: fullscreen, component, video")
	return cmd
}

func run(findRoot rootFinder, args []string, list, live, headed bool, mode string) error {
	if list {
		printList()
		return nil
	}
	targets, err := collectTargets(args, mode)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return fmt.Errorf("no targets specified; use --list to see available options")
	}
	return execSnapshot(findRoot, targets, live, headed)
}
