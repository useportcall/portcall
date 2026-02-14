package subscription

import (
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// StartReset begins a subscription reset workflow.
// For active subscriptions, it marks the status as "resetting".
// For canceled subscriptions, it marks the status as "rollback".
func (s *service) StartReset(input *StartResetInput) (*StartResetResult, error) {
	var sub models.Subscription
	if err := s.db.FindForID(input.SubscriptionID, &sub); err != nil {
		return nil, err
	}

	switch sub.Status {
	case "active":
		return s.startActiveReset(&sub)
	case "canceled":
		return s.startCanceledReset(&sub)
	default:
		return &StartResetResult{Status: "skipped"}, nil
	}
}

func (s *service) startActiveReset(sub *models.Subscription) (*StartResetResult, error) {
	var update models.Subscription
	update.Status = "resetting"
	if err := s.db.Update(&update, "id = ? AND status = ?", sub.ID, "active"); err != nil {
		return nil, err
	}
	return &StartResetResult{
		Subscription: sub,
		Status:       "active",
	}, nil
}

func (s *service) startCanceledReset(sub *models.Subscription) (*StartResetResult, error) {
	if sub.RollbackPlanID == nil {
		id, err := findRollbackPlanID(s.db, sub.AppID, sub.PlanID)
		if err == nil {
			sub.RollbackPlanID = id
		}
	}

	now := time.Now()
	sub.Status = "rollback"
	sub.FinalResetAt = &now
	if err := s.db.Update(sub, "id = ? AND status = ?", sub.ID, "canceled"); err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.FindForID(sub.UserID, &user); err != nil {
		return nil, err
	}

	return &StartResetResult{
		Subscription:   sub,
		User:           &user,
		Status:         "canceled",
		CheckoutURL:    fmt.Sprintf("https://example.com/checkout/%s", sub.PublicID),
		RollbackPlanID: sub.RollbackPlanID,
		FinalResetAt:   now,
	}, nil
}
