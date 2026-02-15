package e2etest

import (
	"testing"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func TestE2E_EntitlementsResetOnImmediatePlanUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	h := NewHarness(t)
	starterID, _ := h.CreatePlanViaAPI("Starter", 3000)
	proID, _ := h.CreatePlanViaAPI("Pro", 8000)
	createMeteredUnitItem(t, h, starterID, 25, 100)
	createMeteredUnitItem(t, h, proID, 90, 1000)
	h.CreateSecretViaAPI()

	starterExtra := h.CreateFeatureViaAPI("starter_extra", false)
	proExtra := h.CreateFeatureViaAPI("pro_extra", false)
	h.CreatePlanFeatureViaAPI(starterExtra, starterID, 1)
	h.CreatePlanFeatureViaAPI(proExtra, proID, 1)

	userID := h.CreateUserViaAPI("Now Upgrade", "upgrade-now@e2e.test")
	h.SetBillingAddressViaAPI(userID)
	sub := createSubscription(t, h, starterID, userID, "pm_upgrade_now")
	processMeterEvent(t, h, userID, "tokens", 12)

	var oldPlan, newPlan models.Plan
	must(t, h.DB.GetForPublicID(h.AppID, starterID, &oldPlan))
	must(t, h.DB.GetForPublicID(h.AppID, proID, &newPlan))
	h.APIPost("/v1/subscriptions/"+sub.PublicID, map[string]any{"plan_id": proID}).
		MustOK(t, "update subscription plan")
	runner := billing.NewRunner(h.DB, h.Crypto, billing.AllStepsWithMeter()...)
	must(t, runner.Run("process_plan_switch", subscription.PlanSwitchInput{
		OldPlanID: oldPlan.ID, NewPlanID: newPlan.ID, SubscriptionID: sub.ID,
	}))

	tokens := h.GetEntitlementViaAPI(userID, "tokens")
	if getFloat(tokens, "usage") != 0 || getFloat(tokens, "quota") != 1000 {
		t.Fatalf("expected tokens usage reset to 0 and quota=1000, got %v", tokens)
	}

	ents := h.ListEntitlementsViaAPI(userID)
	if hasEntitlement(ents, starterExtra) {
		t.Fatalf("starter-only entitlement should be removed on immediate plan change")
	}
	if !hasEntitlement(ents, proExtra) {
		t.Fatalf("pro entitlement missing after immediate plan change")
	}
}

func hasEntitlement(list []map[string]any, id string) bool {
	for _, item := range list {
		if getString(item, "id") == id {
			return true
		}
	}
	return false
}
