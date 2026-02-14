package on_checkout_resolve

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_update"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
	"github.com/useportcall/portcall/libs/go/services/payment"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func resolveSessionHandler(c server.IContext) error {
	var payload checkout_session.ResolvePayload
	if err := json.Unmarshal(c.Payload(), &payload); err != nil {
		return err
	}

	svc := checkout_session.NewService(c.DB(), c.Crypto())
	result, err := svc.Resolve(&payload)
	if err != nil || result.Skipped {
		return err
	}
	return CreatePaymentMethod.Enqueue(c.Queue(), payment.CreateMethodInput{
		AppID:                   result.Session.AppID,
		UserID:                  result.Session.UserID,
		PlanID:                  result.Session.PlanID,
		ExternalPaymentMethodID: result.ExternalPaymentMethodID,
	})
}

func createPaymentMethodHandler(c server.IContext) error {
	var input payment.CreateMethodInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := payment.NewService(c.DB(), c.Crypto())
	result, err := svc.CreatePaymentMethod(&input)
	if err != nil {
		return err
	}
	return UpsertSubscription.Enqueue(c.Queue(), map[string]any{
		"app_id":  result.AppID,
		"user_id": result.UserID,
		"plan_id": result.PlanID,
	})
}

func upsertHandler(c server.IContext) error {
	var input subscription.FindInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := subscription.NewService(c.DB())
	result, err := svc.Find(&input)
	if err != nil {
		return err
	}
	switch result.Action {
	case "create":
		return on_subscription_create.CreateSubscription.Enqueue(c.Queue(), subscription.CreateInput{
			AppID:  result.AppID,
			UserID: result.UserID,
			PlanID: result.PlanID,
		})
	case "update":
		return on_subscription_update.UpdateSubscription.Enqueue(c.Queue(), subscription.UpdateInput{
			SubscriptionID: result.SubscriptionID,
			PlanID:         result.PlanID,
			AppID:          result.AppID,
		})
	}
	return nil
}
