package e2etest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func TestE2E_DiscordNotifications(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	h := NewHarness(t)
	if h.SignupWebhook != nil {
		h.SignupWebhook.Reset()
	}
	if h.BillingWebhook != nil {
		h.BillingWebhook.Reset()
	}

	prepareSignupAccountFlow(t, h)
	h.DashPost("/api/apps", map[string]any{"name": "Discord Signup Project"}).
		MustOK(h.T, "create signup project")

	h.CreateSecretViaAPI()
	apiUserID := h.CreateUserViaAPI("API User", "api-notify@test.dev")
	starterID, _ := h.CreatePlanViaAPI("Discord Starter", 2000)
	proID, _ := h.CreatePlanViaAPI("Discord Pro", 5000)
	h.SetBillingAddressViaAPI(apiUserID)
	csData, _ := h.CreateCheckoutSessionViaAPI(starterID, apiUserID)
	r1 := billing.NewFullRunner(h.DB, h.Crypto)
	must(t, r1.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       getString(csData, "external_session_id"),
		ExternalPaymentMethodID: "pm_notify",
	}))

	var pro models.Plan
	must(t, h.DB.FindFirst(&pro, "public_id = ?", proID))
	var sub models.Subscription
	must(t, h.DB.FindFirst(&sub, "app_id = ? AND status = ?", h.AppID, "active"))
	r2 := billing.NewFullRunner(h.DB, h.Crypto)
	must(t, r2.Run("update_subscription", subscription.UpdateInput{
		SubscriptionID: sub.ID, PlanID: pro.ID, AppID: h.AppID,
	}))

	if h.DiscordMode == discordModeLive {
		t.Logf("live mode: webhook calls sent to configured Discord URLs")
		return
	}
	signup := h.SignupWebhook.WaitForCount(t, 1, 5*time.Second)
	billingMsgs := h.BillingWebhook.WaitForCount(t, 1, 5*time.Second)
	if !hasText(signup, "account signed up") {
		t.Fatalf("signup notifications missing expected text: %v", signup)
	}
	if !hasText(billingMsgs, "upgraded") {
		t.Fatalf("billing notifications missing expected text: %v", billingMsgs)
	}
}

func prepareSignupAccountFlow(t *testing.T, h *Harness) {
	t.Helper()
	var seededApp models.App
	must(t, h.DB.FindFirst(&seededApp, "id = ?", h.AppID))
	var seededAccount models.Account
	must(t, h.DB.FindFirst(&seededAccount, "id = ?", seededApp.AccountID))
	seededAccount.Email = fmt.Sprintf("seeded-%d@portcall.internal", time.Now().UnixNano())
	must(t, h.DB.Save(&seededAccount))
}

func hasText(messages []string, want string) bool {
	for _, m := range messages {
		if strings.Contains(strings.ToLower(m), strings.ToLower(want)) {
			return true
		}
	}
	return false
}
