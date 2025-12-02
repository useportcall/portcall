package handlers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/useportcall/portcall/apps/billing/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type CreateSubscriptionPayload struct {
	PlanID uint `json:"plan_id"`
	UserID uint `json:"user_id"`
}

func CreateSubscription(c server.IContext) error {
	var p CreateSubscriptionPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var user models.User
	if err := c.DB().FindForID(p.UserID, &user); err != nil {
		return err
	}

	if user.BillingAddressID == nil {
		return errors.New("billing address is required")
	}

	var plan models.Plan
	if err := c.DB().FindForID(p.PlanID, &plan); err != nil {
		return err
	}

	err, nextReset := utils.CalculateNextReset(plan.Interval, time.Now())
	if err != nil {
		return err
	}

	subscription := models.Subscription{
		PublicID:             dbx.GenPublicID("sub"),
		AppID:                user.AppID,
		Status:               "active",
		UserID:               user.ID,
		Currency:             plan.Currency,
		BillingAddressID:     *user.BillingAddressID,
		PlanID:               &plan.ID, // TODO: rollback plan?
		BillingInterval:      plan.Interval,
		BillingIntervalCount: plan.IntervalCount,
		InvoiceDueByDays:     plan.InvoiceDueByDays,
		NextResetAt:          nextReset,
	}

	// find free plan to set as rollback
	if err := setRollback(&subscription, &plan, c.DB()); err != nil {
		return err
	}

	if err := c.DB().Create(&subscription); err != nil {
		return err
	}

	if err := c.Queue().Enqueue("create_subscription_items", map[string]any{"subscription_id": subscription.ID, "plan_id": plan.ID}, "billing_queue"); err != nil {
		return err
	}

	if err := c.Queue().Enqueue("create_entitlements", map[string]any{"user_id": user.ID, "plan_id": plan.ID}, "billing_queue"); err != nil {
		return err
	}

	return nil
}

// TODO: improve logic
func setRollback(dest *models.Subscription, src *models.Plan, db dbx.IORM) error {
	if !src.IsFree {
		freePlans := []models.Plan{}
		if err := db.List(&freePlans, "app_id = ? AND is_free = ?", src.AppID, true); err != nil {
			return err
		}

		if len(freePlans) > 0 {
			// Set the first free plan as the rollback plan
			dest.RollbackPlanID = &freePlans[0].ID
		}
	}

	return nil
}
