package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type StartSubscriptionResetPayload struct {
	SubscriptionID uint `json:"subscription_id"`
}

func StartSubscriptionReset(c server.IContext) error {
	var p StartSubscriptionResetPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var subscription models.Subscription
	if err := c.DB().FindForID(p.SubscriptionID, &subscription); err != nil {
		return err
	}

	switch subscription.Status {
	case "active":
		subscription.Status = "resetting"

		if err := c.DB().Update(&subscription, "id = ? AND status = ?", p.SubscriptionID, "active"); err != nil {
			return err
		}

		if err := c.Queue().Enqueue("create_invoice", map[string]any{"subscription_id": subscription.ID}, "billing_queue"); err != nil {
			return nil
		}

		return nil
	case "canceled":
		now := time.Now()
		subscription.Status = "rollback"
		subscription.FinalResetAt = &now

		var user models.User
		if err := c.DB().FindForID(subscription.UserID, &user); err != nil {
			return err
		}

		if err := c.DB().Update(&subscription, "id = ? AND status = ?", p.SubscriptionID, "canceled"); err != nil {
			return err
		}

		if subscription.RollbackPlanID == nil {
			payload := map[string]any{"user_id": subscription.UserID}
			if err := c.Queue().Enqueue("rollback_entitlements", payload, "billing_queue"); err != nil {
				return err
			}
		} else {
			payload := map[string]any{"subscription_id": subscription.ID, "plan_id": *subscription.RollbackPlanID}
			if err := c.Queue().Enqueue("replace_subscription", payload, "billing_queue"); err != nil {
				return err
			}
		}

		payload := map[string]any{
			"checkout_url":  fmt.Sprintf("https://example.com/checkout/%s", subscription.PublicID), //TODO: add logic for resubscribing
			"customer_name": user.Name,
			"year":          2025, // TODO: fix
		}
		if err := c.Queue().Enqueue("subscription_rollback_email", payload, "email_queue"); err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}
