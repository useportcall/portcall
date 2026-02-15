package e2etest

import (
	"os"
	"testing"
)

// TestBrowserHarness starts the browser e2e harness servers.
// Playwright's globalSetup launches: go test -run TestBrowserHarness
// The test prints a JSON config line to stdout and blocks until killed.
func TestBrowserHarness(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping browser harness in short mode")
	}
	if os.Getenv("BROWSER_HARNESS") != "1" {
		t.Skip("skipping browser harness: set BROWSER_HARNESS=1")
	}
	startBrowserHarness(t)
}
