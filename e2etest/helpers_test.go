package e2etest

import (
	"testing"

	"github.com/useportcall/portcall/e2etest/browserharness"
	"github.com/useportcall/portcall/e2etest/harness"
)

// Type aliases so test files reference harness types seamlessly.
type Harness = harness.Harness
type CheckoutFixture = harness.CheckoutFixture
type JSONResponse = harness.JSONResponse

// Function aliases to keep test files concise.
var (
	NewHarness   = harness.NewHarness
	NewTestDB    = harness.NewTestDB
	SeedCheckout = harness.SeedCheckout
)

// HTTP client helper aliases.
var (
	getString  = harness.GetString
	getFloat   = harness.GetFloat
	getSlice   = harness.GetSlice
	mustString = harness.MustString
	dashPath   = harness.DashPath
)

// must fails the test if err is non-nil.
func must(t *testing.T, err error) {
	t.Helper()
	harness.Must(t, err)
}

// Discord mode constants.
const (
	discordModeLive    = harness.DiscordModeLive
	discordModeCapture = harness.DiscordModeCapture
)

// startBrowserHarness launches the browser harness for Playwright.
var startBrowserHarness = browserharness.Start
