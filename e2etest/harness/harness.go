// Package harness provides the core e2e test infrastructure: temporary
// databases, in-process httptest servers, mock dependencies, and HTTP
// client helpers for cross-app integration tests.
package harness

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	apiapp "github.com/useportcall/portcall/apps/api/app"
	checkoutapp "github.com/useportcall/portcall/apps/checkout/app"
	dashapp "github.com/useportcall/portcall/apps/dashboard/app"
	fileapp "github.com/useportcall/portcall/apps/file/app"
	quoteapp "github.com/useportcall/portcall/apps/quote/app"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
)

// TestAESKey is a base64-encoded 32-byte key for test encryption.
var TestAESKey = base64.StdEncoding.EncodeToString(
	[]byte("01234567890123456789012345678901"),
)

// TestAdminAPIKey is the admin bypass key used in tests.
const TestAdminAPIKey = "admin_e2e_test_key"

// Harness starts in-process httptest servers for dashboard, public API,
// checkout, quote, and file. All share the same temp database and crypto.
type Harness struct {
	T                  *testing.T
	DB                 dbx.IORM
	Crypto             cryptox.ICrypto
	Dashboard          *httptest.Server
	API                *httptest.Server
	Checkout           *httptest.Server
	Quote              *httptest.Server
	File               *httptest.Server
	DiscordMode        string
	SignupWebhook      *DiscordCapture
	BillingWebhook     *DiscordCapture
	SnapshotWebhookURL string
	AppPublicID        string // set by SeedBase
	AppID              uint   // set by SeedBase
	APIKey             string // set after CreateSecretViaAPI
}

// Config controls optional harness setup behavior.
type Config struct {
	QuoteTemplates string            // default: "../apps/quote/templates/*.html"
	FileTemplates  string            // default: "../apps/file/templates"
	GinMode        string            // default: "test"
	ExtraEnv       map[string]string // additional env vars to set
}

// NewHarness creates a temp database, starts five httptest servers,
// and returns a harness ready for cross-app e2e tests.
func NewHarness(t *testing.T) *Harness {
	return NewHarnessWithConfig(t, Config{})
}

// NewHarnessWithConfig creates a harness with custom configuration.
func NewHarnessWithConfig(t *testing.T, cfg Config) *Harness {
	t.Helper()

	ginMode := cfg.GinMode
	if ginMode == "" {
		ginMode = "test"
	}
	quoteTpl := cfg.QuoteTemplates
	if quoteTpl == "" {
		quoteTpl = "../apps/quote/templates/*.html"
	}
	fileTpl := cfg.FileTemplates
	if fileTpl == "" {
		fileTpl = "../apps/file/templates"
	}

	t.Setenv("AES_ENCRYPTION_KEY", TestAESKey)
	t.Setenv("ADMIN_API_KEY", TestAdminAPIKey)
	t.Setenv("INVOICE_APP_URL", "https://e2e.example.com")
	t.Setenv("GIN_MODE", ginMode)
	t.Setenv("E2E_MODE", "true")
	for k, v := range cfg.ExtraEnv {
		t.Setenv(k, v)
	}

	jwks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"keys":[]}`))
	}))
	t.Cleanup(jwks.Close)
	t.Setenv("KEYCLOAK_API_URL", jwks.URL)

	db := NewTestDB(t)
	crypto, err := cryptox.New()
	if err != nil {
		t.Fatalf("init crypto: %v", err)
	}
	q := &NopQueue{}
	store := NewMemoryStore()

	dashRouter, err := dashapp.NewRouter(db, crypto, q)
	if err != nil {
		t.Fatalf("init dashboard router: %v", err)
	}
	dashRouter.SetStore(store)
	dash := httptest.NewServer(dashRouter.Handler())
	t.Cleanup(dash.Close)

	checkoutSrv := httptest.NewServer(checkoutapp.NewRouter(db, crypto, q).Handler())
	t.Cleanup(checkoutSrv.Close)
	t.Setenv("CHECKOUT_URL", checkoutSrv.URL)
	t.Setenv("CHECKOUT_APP_URL", checkoutSrv.URL)

	quoteSrv := httptest.NewServer(quoteapp.NewRouter(db, crypto, q, store, quoteTpl).Handler())
	t.Cleanup(quoteSrv.Close)
	t.Setenv("QUOTE_APP_URL", quoteSrv.URL)

	fileSrv := httptest.NewServer(fileapp.NewRouter(db, store, fileTpl).Handler())
	t.Cleanup(fileSrv.Close)
	t.Setenv("INVOICE_APP_URL", fileSrv.URL)

	api := httptest.NewServer(apiapp.NewRouter(db, crypto, q).Handler())
	t.Cleanup(api.Close)

	h := &Harness{
		T: t, DB: db, Crypto: crypto,
		Dashboard: dash, API: api, Checkout: checkoutSrv,
		Quote: quoteSrv, File: fileSrv,
	}
	SetupDiscord(t, h)
	h.SeedBase()
	return h
}
