package e2etest

import (
	"testing"
	"time"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

// TestE2E_PlanUpgradeFlow exercises the full plan-upgrade lifecycle via
// real HTTP calls: Dashboard → Public API → Billing (initial) → Billing (upgrade)
func TestE2E_PlanUpgradeFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	h := NewHarness(t)

	allSagas := billing.AllSteps()

	// ── STAGE 1: Dashboard API — create two plans ──
	starterID, _ := h.CreatePlanViaAPI("Starter", 4900)
	proID, _ := h.CreatePlanViaAPI("Pro", 9900)
	t.Logf("✓ STAGE 1 — Dashboard API: Starter=%s ($49) Pro=%s ($99)", starterID, proID)

	// ── STAGE 2: Create secret + user with billing address ──
	h.CreateSecretViaAPI()
	userID := h.CreateUserViaAPI("Jordan", "jordan@upgrade.ai")
	h.SetBillingAddressViaAPI(userID)
	t.Logf("✓ STAGE 2 — created secret + user %s with billing address", userID)

	// ── STAGE 3: Public API — create checkout session for Starter ──
	csData, _ := h.CreateCheckoutSessionViaAPI(starterID, userID)
	t.Logf("✓ STAGE 3 — Public API: checkout session %s (Starter)", getString(csData, "id"))

	// ── STAGE 4: Billing — resolve checkout → initial subscription ──
	r1 := billing.NewRunner(h.DB, h.Crypto, allSagas...)
	err := r1.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       getString(csData, "external_session_id"),
		ExternalPaymentMethodID: "pm_upg_001",
	})
	if err != nil {
		t.Fatalf("Billing: initial saga: %v", err)
	}
	var sub models.Subscription
	must(t, h.DB.FindFirst(&sub, "app_id = ?", h.AppID))
	t.Logf("✓ STAGE 4 — Billing: subscription created (id=%d, plan=Starter)", sub.ID)

	// Set billing cycle for proration calculations.
	setSubscriptionTiming(t, h.DB, sub.ID)

	// Look up the Pro plan's internal ID for the upgrade saga.
	var proPlan models.Plan
	must(t, h.DB.FindFirst(&proPlan, "public_id = ?", proID))

	// ── STAGE 5: Billing — upgrade Starter → Pro ──
	r2 := billing.NewRunner(h.DB, h.Crypto, allSagas...)
	err = r2.Run("update_subscription", subscription.UpdateInput{
		SubscriptionID: sub.ID, PlanID: proPlan.ID, AppID: h.AppID,
	})
	if err != nil {
		t.Fatalf("Billing: upgrade saga: %v", err)
	}
	for _, name := range []string{"update_subscription", "pay_invoice"} {
		if !r2.HasTask(name) {
			t.Errorf("Billing: missing upgrade task %q; ran: %v", name, r2.Executed)
		}
	}
	if len(r2.TaskPayloads("pay_invoice")) < 1 {
		t.Errorf("Billing: expected at least one pay_invoice task; ran: %v", r2.Executed)
	}
	t.Logf("✓ STAGE 5 — Billing: upgrade saga (%d tasks)", len(r2.Executed))

	// ── STAGE 6: Verify final state ──
	verifyUpgradeFinalState(t, h, sub.ID, proPlan.ID)
}

func setSubscriptionTiming(t *testing.T, db interface{ Exec(string, ...any) error }, id uint) {
	t.Helper()
	now := time.Now()
	must(t, db.Exec(
		"UPDATE subscriptions SET last_reset_at = ?, next_reset_at = ? WHERE id = ?",
		now.AddDate(0, 0, -15), now.AddDate(0, 0, 15), id,
	))
}

func verifyUpgradeFinalState(t *testing.T, h *Harness, subID, proID uint) {
	t.Helper()
	var s models.Subscription
	must(t, h.DB.FindForID(subID, &s))
	if *s.PlanID != proID {
		t.Fatalf("Final: plan should be Pro (%d), got %d", proID, *s.PlanID)
	}
	var invoices []models.Invoice
	must(t, h.DB.List(&invoices, "subscription_id = ?", subID))
	if len(invoices) < 2 {
		t.Fatalf("Final: expected ≥2 invoices, got %d", len(invoices))
	}
	for _, inv := range invoices {
		if inv.Status != "paid" {
			t.Errorf("Final: invoice %d status=%s, want paid", inv.ID, inv.Status)
		}
	}
	t.Logf("✓ STAGE 6 — Final: plan=Pro, invoices=%d (all paid)", len(invoices))
}
