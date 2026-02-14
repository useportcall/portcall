package on_subscription_reset

import (
	"encoding/json"
	"time"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func findHandler(c server.IContext) error {
	svc := subscription.NewService(c.DB())
	result, err := svc.FindResets()
	if err != nil {
		return err
	}
	for _, id := range result.SubscriptionIDs {
		if err := StartReset.Enqueue(c.Queue(), subscription.StartResetInput{
			SubscriptionID: id,
		}); err != nil {
			return err
		}
	}
	return nil
}

func startResetHandler(c server.IContext) error {
	var input subscription.StartResetInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := subscription.NewService(c.DB())
	result, err := svc.StartReset(&input)
	if err != nil {
		return err
	}

	switch result.Status {
	case "active":
		if err := on_invoice_create.CreateInvoice.Enqueue(c.Queue(), map[string]any{
			"subscription_id": result.Subscription.ID,
		}); err != nil {
			return err
		}
		return EndReset.Enqueue(c.Queue(), subscription.EndResetInput{
			SubscriptionID: result.Subscription.ID,
		})
	case "canceled":
		if result.RollbackPlanID != nil {
			if err := on_subscription_create.CreateSubscription.Enqueue(c.Queue(), subscription.CreateInput{
				AppID:  result.Subscription.AppID,
				UserID: result.Subscription.UserID,
				PlanID: *result.RollbackPlanID,
			}); err != nil {
				return err
			}
		}
		return rollbackEmail.Enqueue(c.Queue(), map[string]any{
			"checkout_url":  result.CheckoutURL,
			"customer_name": result.User.Name,
			"year":          time.Now().Year(),
		})
	}
	return nil
}

func endResetHandler(c server.IContext) error {
	var input subscription.EndResetInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := subscription.NewService(c.DB())
	result, err := svc.EndReset(&input)
	if err != nil {
		return err
	}
	return syncEntitlements(c, result.Subscription.UserID, result.AppliedPlanID)
}
