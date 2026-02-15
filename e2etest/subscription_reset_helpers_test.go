package e2etest

import (
	"testing"
	"time"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
)

type resetFixture struct {
	ProPlanID      uint
	SubOneID       uint
	SubTwoID       uint
	SubTwoLast     time.Time
	SubTwoNext     time.Time
	MeteredFeature string
	UserOnePublic  string
}

func buildResetFixture(t *testing.T, h *Harness) *resetFixture {
	t.Helper()
	starterID, _ := h.CreatePlanViaAPI("Starter", 3000)
	proID, _ := h.CreatePlanViaAPI("Pro", 7000)
	createMeteredUnitItem(t, h, starterID, 25, 100)
	createMeteredUnitItem(t, h, proID, 75, 500)
	h.CreateSecretViaAPI()

	userOne := h.CreateUserViaAPI("Reset One", "reset-one@e2e.test")
	userTwo := h.CreateUserViaAPI("Reset Two", "reset-two@e2e.test")
	h.SetBillingAddressViaAPI(userOne)
	h.SetBillingAddressViaAPI(userTwo)

	subOne := createSubscription(t, h, starterID, userOne, "pm_reset_001")
	subTwo := createSubscription(t, h, starterID, userTwo, "pm_reset_002")
	h.APIPost("/v1/subscriptions/"+subOne.PublicID, map[string]any{
		"plan_id":             proID,
		"apply_at_next_reset": true,
	}).MustOK(t, "schedule plan switch")

	processMeterEvent(t, h, userOne, "tokens", 12)
	now := time.Now()
	must(t, h.DB.Exec("UPDATE subscriptions SET next_reset_at=?, last_reset_at=? WHERE id=?",
		now.Add(-time.Minute), now.AddDate(0, 0, -30), subOne.ID))
	must(t, h.DB.Exec("UPDATE subscriptions SET next_reset_at=?, last_reset_at=? WHERE id=?",
		now.Add(24*time.Hour), now.AddDate(0, 0, -3), subTwo.ID))

	var proPlan models.Plan
	must(t, h.DB.GetForPublicID(h.AppID, proID, &proPlan))
	var subTwoFresh models.Subscription
	must(t, h.DB.FindForID(subTwo.ID, &subTwoFresh))
	return &resetFixture{
		ProPlanID:      proPlan.ID,
		SubOneID:       subOne.ID,
		SubTwoID:       subTwo.ID,
		SubTwoLast:     subTwoFresh.LastResetAt,
		SubTwoNext:     subTwoFresh.NextResetAt,
		MeteredFeature: "tokens",
		UserOnePublic:  userOne,
	}
}

func createMeteredUnitItem(t *testing.T, h *Harness, planID string, unitAmount, quota int64) {
	t.Helper()
	h.DashPost(dashPath(h.AppPublicID, "plan-items"), map[string]any{
		"plan_id":       planID,
		"pricing_model": "unit",
		"unit_amount":   unitAmount,
		"public_title":  "API Calls",
		"quota":         quota,
	}).MustOK(t, "create metered plan item")
}

func createSubscription(t *testing.T, h *Harness, planID, userID, pmID string) models.Subscription {
	t.Helper()
	cs, _ := h.CreateCheckoutSessionViaAPI(planID, userID)
	r := billing.NewRunner(h.DB, h.Crypto, billing.AllStepsWithMeter()...)
	must(t, r.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       getString(cs, "external_session_id"),
		ExternalPaymentMethodID: pmID,
	}))
	var user models.User
	must(t, h.DB.GetForPublicID(h.AppID, userID, &user))
	var sub models.Subscription
	must(t, h.DB.FindFirst(&sub, "app_id = ? AND user_id = ? AND status = ?", h.AppID, user.ID, "active"))
	return sub
}
