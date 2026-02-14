package on_subscription_create

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func createSubscriptionHandler(c server.IContext) error {
	var input subscription.CreateInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	return subscription.RunCreateFlow(&subscription.BillingFlowInput{
		DB:           c.DB(),
		Crypto:       c.Crypto(),
		Queue:        c.Queue(),
		AppID:        input.AppID,
		UserID:       input.UserID,
		PlanID:       input.PlanID,
		AsyncPayment: true,
	})
}

func createItemsHandler(c server.IContext) error {
	var input subscription.CreateItemsInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := subscription.NewService(c.DB())
	result, err := svc.CreateItems(&input)
	if err != nil {
		return err
	}
	if result.ItemCount == 0 {
		return nil
	}
	return on_invoice_create.CreateInvoice.Enqueue(c.Queue(), map[string]any{
		"subscription_id": result.SubscriptionID,
	})
}

func createEntitlementsHandler(c server.IContext) error {
	var input entitlement.StartUpsertInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := entitlement.NewService(c.DB())
	result, err := svc.StartUpsert(&input)
	if err != nil || !result.HasMore {
		return err
	}
	return CreateSingleEntitlement.Enqueue(c.Queue(), entitlement.UpsertInput{
		UserID: result.UserID,
		Index:  0,
		Values: result.PlanFeatureIDs,
	})
}

func createSingleEntitlementHandler(c server.IContext) error {
	var input entitlement.UpsertInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := entitlement.NewService(c.DB())
	result, err := svc.Upsert(&input)
	if err != nil || !result.HasMore {
		return err
	}
	return CreateSingleEntitlement.Enqueue(c.Queue(), entitlement.UpsertInput{
		UserID: result.UserID,
		Index:  result.NextIndex,
		Values: result.PlanFeatureIDs,
	})
}
