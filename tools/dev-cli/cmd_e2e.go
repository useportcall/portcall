package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newE2ECmd() *cobra.Command {
	return &cobra.Command{
		Use:   "e2e [filter]",
		Short: "Run cross-app e2e tests against a temporary database",
		Long: `Run end-to-end tests. A temporary Postgres database is created
per test, migrations run, tests execute, and the database is dropped.

The portcall postgres_instance container is started automatically if needed.

Examples:
  dev-cli e2e              # all e2e tests
  dev-cli e2e flow         # only tests matching "Flow"
  dev-cli e2e upgrade      # only tests matching "Upgrade"`,
		Args: cobra.MaximumNArgs(1),
		RunE: e2eCommand,
	}
}

func e2eCommand(cmd *cobra.Command, args []string) error {
	pkgs := e2eTargets["all"]

	root, err := findRootDir()
	if err != nil {
		return err
	}

	fmt.Println("Running cross-app e2e tests...")
	runPattern := "TestE2E"
	if len(args) > 0 {
		// Allow partial match, e.g. "flow" → "TestE2E.*(?i)flow"
		runPattern = fmt.Sprintf("TestE2E.*(?i)%s", args[0])
		fmt.Printf("Filter: %s\n", runPattern)
	}
	goArgs := append([]string{"test", "-v", "-count=1", "-timeout", "120s", "-run", runPattern}, pkgs...)
	c := exec.Command("go", goArgs...)
	c.Dir = root
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	if err := c.Run(); err != nil {
		return fmt.Errorf("e2e tests failed: %w", err)
	}
	fmt.Println("✅ E2E tests passed!")
	return nil
}
