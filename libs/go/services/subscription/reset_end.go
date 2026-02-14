package subscription

import (
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// EndReset completes a subscription reset by updating the last
// and next reset times and restoring the active status.
func (s *service) EndReset(input *EndResetInput) (*EndResetResult, error) {
	var sub models.Subscription
	if err := s.db.FindForID(input.SubscriptionID, &sub); err != nil {
		return nil, err
	}

	now := time.Now()
	var (
		nextReset     *time.Time
		appliedPlanID *uint
	)
	if err := s.db.Txn(func(tx dbx.IORM) error {
		planID, err := applyScheduledPlan(tx, &sub)
		if err != nil {
			return err
		}
		appliedPlanID = planID

		nextReset, err = NextReset(sub.CreatedAt, sub.BillingInterval, now)
		if err != nil {
			return fmt.Errorf("failed to calculate next reset: %w", err)
		}

		sub.Status = "active"
		sub.LastResetAt = now
		sub.NextResetAt = *nextReset
		if err := tx.Save(&sub); err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}

		return resetUsageRecords(tx, sub.ID, now, *nextReset)
	}); err != nil {
		return nil, err
	}

	return &EndResetResult{
		Subscription:  &sub,
		LastResetAt:   now,
		NextResetAt:   *nextReset,
		AppliedPlanID: appliedPlanID,
	}, nil
}
