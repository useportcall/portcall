package e2etest

import (
	"testing"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
)

// TestE2E_FullCheckoutFlow exercises the cross-service checkout flow via
// real HTTP calls: Dashboard API → Public API → Checkout API → Billing Sagas
func TestE2E_FullCheckoutFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	h := NewHarness(t)

	// ── STAGE 1: Dashboard API — create & publish plan ($49/mo) ──
	planID, _ := h.CreatePlanViaAPI("AI Starter", 4900)
	t.Logf("✓ STAGE 1 — Dashboard API: created plan %q ($49/mo)", planID)

	// ── STAGE 2: Dashboard API — create API secret ──
	key := h.CreateSecretViaAPI()
	t.Logf("✓ STAGE 2 — Dashboard API: created API secret (len=%d)", len(key))

	// ── STAGE 3: Public API — create user + billing address ──
	userID := h.CreateUserViaAPI("Alex", "alex@ai.app")
	h.SetBillingAddressViaAPI(userID)
	t.Logf("✓ STAGE 3 — Public API: created user %q with billing address", userID)

	// ── STAGE 4: Public API — create checkout session ──
	csData, csToken := h.CreateCheckoutSessionViaAPI(planID, userID)
	csID := mustString(t, csData, "id")
	provider := getString(csData, "external_provider")
	if provider != "local" {
		t.Fatalf("expected provider local, got %q", provider)
	}
	t.Logf("✓ STAGE 4 — Public API: created checkout session %q (provider=%s)", csID, provider)

	// ── STAGE 5: Checkout API — fetch session by public ID ──
	csGet := h.CheckoutGet("/api/checkout-sessions/"+csID, csToken)
	fetched := csGet.MustOK(t, "get checkout session")
	if getString(fetched, "id") != csID {
		t.Fatalf("Checkout API: session ID mismatch")
	}
	t.Logf("✓ STAGE 5 — Checkout API: fetched session %q", csID)

	// ── STAGE 6: Billing — run checkout-resolve saga chain ──
	extSessionID := getString(csData, "external_session_id")
	runner := billing.NewFullRunner(h.DB, h.Crypto)
	err := runner.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       extSessionID,
		ExternalPaymentMethodID: "pm_flow_001",
	})
	if err != nil {
		t.Fatalf("Billing: saga chain failed: %v", err)
	}
	expect := []string{
		"resolve_checkout_session", "create_payment_method", "upsert_subscription",
		"create_subscription", "pay_invoice",
	}
	for _, name := range expect {
		if !runner.HasTask(name) {
			t.Errorf("Billing: missing task %q; ran: %v", name, runner.Executed)
		}
	}
	t.Logf("✓ STAGE 6 — Billing: saga completed (%d tasks)", len(runner.Executed))

	// ── STAGE 7: Verify final state ──
	verifyFlowFinalState(t, h)
}

func verifyFlowFinalState(t *testing.T, h *Harness) {
	t.Helper()
	var sub models.Subscription
	must(t, h.DB.FindFirst(&sub, "app_id = ?", h.AppID))
	if sub.Status != "active" {
		t.Fatalf("Final: subscription status %q, want active", sub.Status)
	}
	var inv models.Invoice
	must(t, h.DB.FindFirst(&inv, "subscription_id = ?", sub.ID))
	if inv.Status != "paid" {
		t.Fatalf("Final: invoice status %q, want paid", inv.Status)
	}
	t.Logf("✓ STAGE 7 — Final: subscription=%s invoice=%s", sub.Status, inv.Status)
}
