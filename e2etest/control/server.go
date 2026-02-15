// Package control provides an HTTP control server for the browser
// harness. Playwright tests call these endpoints to trigger server-side
// actions that cannot be done from the browser (e.g. billing sagas).
package control

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/e2etest/harness"
	checkout_session "github.com/useportcall/portcall/libs/go/services/checkout_session"
)

// Server wraps the harness in an HTTP control endpoint.
type Server struct {
	H *harness.Harness
}

type resolveCheckoutRequest struct {
	ExternalSessionID string `json:"external_session_id"`
}

// NewServer starts an httptest server with e2e-only control endpoints.
func NewServer(h *harness.Harness) *httptest.Server {
	s := &Server{H: h}
	mux := http.NewServeMux()
	mux.HandleFunc("/e2e/resolve-checkout", s.handleResolve)
	mux.HandleFunc("/e2e/upgrade-subscription", s.handleUpgradeSubscription)
	mux.HandleFunc("/e2e/discord/mode", s.handleDiscordMode)
	mux.HandleFunc("/e2e/discord/messages", s.handleDiscordMessages)
	mux.HandleFunc("/e2e/discord/reset", s.handleDiscordReset)
	mux.HandleFunc("/e2e/signup/prepare", s.handlePrepareSignup)
	mux.HandleFunc("/e2e/user/payment-status", s.handleSetUserPaymentStatus)
	mux.HandleFunc("/e2e/snapshot/config", s.handleSnapshotConfig)
	mux.HandleFunc("/e2e/checkout-session", s.handleCheckoutSession)
	mux.HandleFunc("/e2e/seed-invoice", s.handleSeedInvoice)
	return httptest.NewServer(mux)
}

func (s *Server) handleResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req resolveCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	runner := billing.NewFullRunner(s.H.DB, s.H.Crypto)
	err := runner.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       req.ExternalSessionID,
		ExternalPaymentMethodID: "pm_browser_e2e",
	})
	if err != nil {
		http.Error(w, "saga failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{"ok": true, "tasks": runner.Executed}
	json.NewEncoder(w).Encode(resp)
}
