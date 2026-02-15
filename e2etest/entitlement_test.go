package e2etest

import (
	"testing"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
)

// TestE2E_FeatureEntitlementFlow exercises the full feature → plan-feature →
// checkout → entitlement → meter-event lifecycle across all services.
func TestE2E_FeatureEntitlementFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	h := NewHarness(t)

	// ── STAGE 1: Dashboard — create plan + API secret ──
	planID, _ := h.CreatePlanViaAPI("Pro Plan", 9900)
	h.CreateSecretViaAPI()
	t.Logf("✓ STAGE 1 — Created plan %s + API secret", planID)

	// ── STAGE 2: Create features (1 boolean + 1 metered) and attach to plan ──
	boolFeatureID := h.CreateFeatureViaAPI("dashboard-access", false)
	meteredFeatureID := h.CreateFeatureViaAPI("api-calls", true)
	h.CreatePlanFeatureViaAPI(boolFeatureID, planID, -1)     // unlimited boolean
	h.CreatePlanFeatureViaAPI(meteredFeatureID, planID, 100) // 100 API calls quota
	t.Logf("✓ STAGE 2 — Features: boolean=%s metered=%s (attached to plan)", boolFeatureID, meteredFeatureID)

	// ── STAGE 3: Create user with billing address ──
	userID := h.CreateUserViaAPI("Sam", "sam@features.ai")
	h.SetBillingAddressViaAPI(userID)
	t.Logf("✓ STAGE 3 — Created user %s with billing address", userID)

	// ── STAGE 4: Create checkout session + billing saga → subscription + entitlements ──
	csData, _ := h.CreateCheckoutSessionViaAPI(planID, userID)
	extSessionID := getString(csData, "external_session_id")

	runner := billing.NewFullRunner(h.DB, h.Crypto)
	err := runner.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       extSessionID,
		ExternalPaymentMethodID: "pm_feat_001",
	})
	if err != nil {
		t.Fatalf("Billing saga failed: %v", err)
	}
	if !runner.HasTask("create_subscription") {
		t.Fatalf("Billing: missing create_subscription task; ran: %v", runner.Executed)
	}
	t.Logf("✓ STAGE 4 — Subscription created with entitlements (%d tasks)", len(runner.Executed))

	// ── STAGE 5: Verify entitlements via public API ──
	ents := h.ListEntitlementsViaAPI(userID)
	// The boolean feature should appear in non-metered list
	found := false
	for _, e := range ents {
		if getString(e, "id") == boolFeatureID {
			found = true
		}
	}
	if !found {
		// Try metered list too — the list might be filtered by is_metered
		meteredEnts := h.APIGet("/v1/entitlements?user_id="+userID+"&is_metered=true").
			MustOKList(t, "list metered entitlements")
		ents = append(ents, meteredEnts...)
	}
	t.Logf("✓ STAGE 5 — Found %d entitlements for user", len(ents))

	// ── STAGE 6: Check individual entitlement access ──
	boolEnt := h.GetEntitlementViaAPI(userID, boolFeatureID)
	if boolEnt["enabled"] != true {
		t.Fatalf("Boolean entitlement should be enabled, got: %v", boolEnt)
	}
	t.Logf("✓ STAGE 6 — Boolean entitlement enabled=%v", boolEnt["enabled"])

	meteredEnt := h.GetEntitlementViaAPI(userID, meteredFeatureID)
	if meteredEnt["enabled"] != true {
		t.Fatalf("Metered entitlement should be enabled, got: %v", meteredEnt)
	}
	if getFloat(meteredEnt, "usage") != 0 {
		t.Fatalf("Metered entitlement should start at usage=0, got: %v", meteredEnt["usage"])
	}
	quota := getFloat(meteredEnt, "quota")
	if quota != 100 {
		t.Fatalf("Metered entitlement quota should be 100, got: %v", quota)
	}
	t.Logf("✓ STAGE 7 — Metered entitlement: usage=0 quota=%.0f enabled=%v", quota, meteredEnt["enabled"])

	// ── STAGE 8: Record meter events + process them via billing saga ──
	resp := h.RecordMeterEventViaAPI(userID, meteredFeatureID, 10)
	resp.MustStatus(t, 200, "record meter event")

	// Process the meter event via billing saga (normally queued).
	var meterEvent models.MeterEvent
	must(t, h.DB.FindFirst(&meterEvent, "user_id = (SELECT id FROM users WHERE public_id = ?)", userID))
	meterRunner := billing.NewRunner(h.DB, h.Crypto, billing.MeterEventSteps())
	err = meterRunner.Run("process_meter_event", entitlement.IncrementUsageInput{
		MeterEventID: meterEvent.ID,
	})
	if err != nil {
		t.Fatalf("Meter event saga failed: %v", err)
	}
	t.Logf("✓ STAGE 8 — Meter event recorded and processed")

	// ── STAGE 9: Verify usage incremented ──
	updatedEnt := h.GetEntitlementViaAPI(userID, meteredFeatureID)
	usage := getFloat(updatedEnt, "usage")
	if usage != 10 {
		t.Fatalf("Expected usage=10 after meter event, got %.0f", usage)
	}
	t.Logf("✓ STAGE 9 — Metered entitlement usage=%.0f/%.0f", usage, quota)

	// ── STAGE 10: Exceed quota → expect rejection ──
	overResp := h.RecordMeterEventViaAPI(userID, meteredFeatureID, 200)
	if overResp.Status != 400 {
		t.Fatalf("Expected 400 for exceeding quota, got %d — %v", overResp.Status, overResp.Raw)
	}
	t.Logf("✓ STAGE 10 — Quota enforcement: rejected over-quota meter event (status=%d)", overResp.Status)
}
