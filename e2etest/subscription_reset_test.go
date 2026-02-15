package e2etest

import (
	"testing"
	"time"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func TestE2E_SubscriptionResetFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	h := NewHarness(t)
	f := buildResetFixture(t, h)

	var beforeItem models.SubscriptionItem
	must(t, h.DB.FindFirst(&beforeItem, "subscription_id = ? AND pricing_model = ?", f.SubOneID, "unit"))
	if beforeItem.Usage == 0 {
		t.Fatalf("expected pre-reset metered usage > 0")
	}
	var beforeMeter models.BillingMeter
	must(t, h.DB.FindFirst(&beforeMeter, "subscription_id = ?", f.SubOneID))
	if beforeMeter.Usage == 0 {
		t.Fatalf("expected pre-reset billing meter usage > 0")
	}

	runner := billing.NewRunner(h.DB, h.Crypto, billing.AllStepsWithMeter()...)
	must(t, runner.Run("find_subscriptions_to_reset", map[string]any{}))
	if got := len(runner.TaskPayloads("start_subscription_reset")); got != 1 {
		t.Fatalf("expected exactly 1 subscription reset, got %d", got)
	}

	var subOne models.Subscription
	must(t, h.DB.FindForID(f.SubOneID, &subOne))
	if subOne.PlanID == nil || *subOne.PlanID != f.ProPlanID || subOne.ScheduledPlanID != nil {
		t.Fatalf("scheduled plan not applied at reset")
	}
	if subOne.LastResetAt.Before(time.Now().Add(-time.Minute)) || !subOne.NextResetAt.After(time.Now()) {
		t.Fatalf("subscription reset timestamps not updated correctly")
	}

	var ent models.Entitlement
	must(t, h.DB.FindFirst(&ent, "user_id = (SELECT id FROM users WHERE public_id = ?) AND feature_public_id = ?",
		f.UserOnePublic, f.MeteredFeature))
	if ent.Usage != 0 || ent.Quota != 500 {
		t.Fatalf("entitlement reset/upsert mismatch: usage=%d quota=%d", ent.Usage, ent.Quota)
	}

	var afterItem models.SubscriptionItem
	must(t, h.DB.FindFirst(&afterItem, "subscription_id = ? AND pricing_model = ?", f.SubOneID, "unit"))
	if afterItem.Usage != 0 {
		t.Fatalf("subscription item usage should reset to 0, got %d", afterItem.Usage)
	}
	var afterMeter models.BillingMeter
	must(t, h.DB.FindFirst(&afterMeter, "subscription_id = ?", f.SubOneID))
	if afterMeter.Usage != 0 {
		t.Fatalf("billing meter usage should reset to 0, got %d", afterMeter.Usage)
	}

	var subTwo models.Subscription
	must(t, h.DB.FindForID(f.SubTwoID, &subTwo))
	if !subTwo.LastResetAt.Equal(f.SubTwoLast) || !subTwo.NextResetAt.Equal(f.SubTwoNext) {
		t.Fatalf("non-due subscription should not be reset")
	}
}
