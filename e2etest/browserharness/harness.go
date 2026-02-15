package browserharness

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/useportcall/portcall/e2etest/control"
	"github.com/useportcall/portcall/e2etest/harness"
)

// Config is the JSON config emitted for Playwright tests to consume.
type Config struct {
	DashboardURL string `json:"dashboard_url"`
	APIURL       string `json:"api_url"`
	CheckoutURL  string `json:"checkout_url"`
	QuoteURL     string `json:"quote_url"`
	FileURL      string `json:"file_url"`
	ControlURL   string `json:"control_url"`
	AdminAPIKey  string `json:"admin_api_key"`
	AppID        uint   `json:"app_id"`
	AppPublicID  string `json:"app_public_id"`
}

// Start launches all servers, prints a JSON config line, and blocks until SIGTERM.
func Start(t *testing.T) {
	root := harness.FindRootDir()

	h := harness.NewHarnessWithConfig(t, harness.Config{
		QuoteTemplates: filepath.Join(root, "apps/quote/templates/*.html"),
		FileTemplates:  filepath.Join(root, "apps/file/templates"),
		GinMode:        "release",
		ExtraEnv: map[string]string{
			"DASHBOARD_STATIC_DIR": filepath.Join(root, "apps/dashboard/frontend/dist"),
			"CHECKOUT_STATIC_DIR":  filepath.Join(root, "apps/checkout/frontend/dist"),
		},
	})

	ctrl := control.NewServer(h)
	t.Cleanup(ctrl.Close)
	t.Setenv("E2E_APP_ID", fmt.Sprintf("%d", h.AppID))

	out, _ := json.Marshal(Config{
		DashboardURL: h.Dashboard.URL,
		APIURL:       h.API.URL,
		CheckoutURL:  h.Checkout.URL,
		QuoteURL:     h.Quote.URL,
		FileURL:      h.File.URL,
		ControlURL:   ctrl.URL,
		AdminAPIKey:  harness.TestAdminAPIKey,
		AppID:        h.AppID,
		AppPublicID:  h.AppPublicID,
	})
	fmt.Fprintln(os.Stdout, string(out))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}
