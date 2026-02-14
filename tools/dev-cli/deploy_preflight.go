package main

import "fmt"

func runPreflightSuites() error {
	info("Running pre-deploy test suites...")
	if deployOpts.runUnitTests {
		if err := runUnitSuite(); err != nil {
			return err
		}
	}
	if deployOpts.runIntegTests {
		if err := runSuite("integration", "go", "test", "-tags=integration", "./libs/go/services/..."); err != nil {
			return err
		}
	}
	if deployOpts.runE2ETests {
		if err := runSuite("e2e", "make", "e2e"); err != nil {
			return err
		}
	}
	ok("Pre-deploy suites passed")
	return nil
}

func runSuite(name string, cmd string, args ...string) error {
	info("Running %s tests...", name)
	if err := runCmd(cmd, args...); err != nil {
		return fmt.Errorf("%s tests failed: %w", name, err)
	}
	ok("%s tests passed", name)
	return nil
}

func runUnitSuite() error {
	info("Running unit tests...")
	script := `set -euo pipefail
mods=$(awk '/^\t\.\//{gsub(/^\t/,""); print}' go.work | grep -v '^./e2etest$')
for m in $mods; do
  echo "==> $m"
  (cd "$m" && go test ./...)
done`
	if err := runShell(script); err != nil {
		return fmt.Errorf("unit tests failed: %w", err)
	}
	ok("unit tests passed")
	return nil
}
