package control

import (
	"encoding/json"
	"net/http"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

type upgradeSubscriptionRequest struct {
	UserID string `json:"user_id"`
	PlanID string `json:"plan_id"`
}

func (s *Server) handleUpgradeSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req upgradeSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := s.H.DB.FindFirst(&user, "app_id = ? AND public_id = ?", s.H.AppID, req.UserID); err != nil {
		http.Error(w, "user not found: "+err.Error(), http.StatusNotFound)
		return
	}
	var plan models.Plan
	if err := s.H.DB.FindFirst(&plan, "app_id = ? AND public_id = ?", s.H.AppID, req.PlanID); err != nil {
		http.Error(w, "plan not found: "+err.Error(), http.StatusNotFound)
		return
	}
	var sub models.Subscription
	if err := s.H.DB.FindFirst(&sub, "app_id = ? AND user_id = ? AND status = ?", s.H.AppID, user.ID, "active"); err != nil {
		http.Error(w, "subscription not found: "+err.Error(), http.StatusNotFound)
		return
	}

	runner := billing.NewFullRunner(s.H.DB, s.H.Crypto)
	err := runner.Run("update_subscription", subscription.UpdateInput{
		SubscriptionID: sub.ID,
		PlanID:         plan.ID,
		AppID:          s.H.AppID,
	})
	if err != nil {
		http.Error(w, "upgrade failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "tasks": runner.Executed})
}
