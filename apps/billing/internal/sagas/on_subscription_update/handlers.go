package on_subscription_update

import (
	"encoding/json"

	"fmt"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/libs/go/discordx"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
	"github.com/useportcall/portcall/libs/go/services/invoice"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func updateHandler(c server.IContext) error {
	var input subscription.UpdateInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	return subscription.RunUpdateFlow(&subscription.BillingFlowInput{
		DB:           c.DB(),
		Crypto:       c.Crypto(),
		Queue:        c.Queue(),
		AppID:        input.AppID,
		PlanID:       input.PlanID,
		AsyncPayment: true,
		NotifyUpgrade: func(userID, newPlanID uint) {
			discordx.SendFromEnvAsync(
				"DISCORD_WEBHOOK_URL_BILLING",
				fmt.Sprintf("💰 User %d upgraded to plan %d", userID, newPlanID),
			)
		},
	}, input.SubscriptionID)
}

func planSwitchHandler(c server.IContext) error {
	var input subscription.PlanSwitchInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := subscription.NewService(c.DB())
	result, err := svc.PlanSwitch(&input)
	if err != nil {
		return err
	}

	if result.IsUpgrade {
		discordx.SendFromEnvAsync(
			"DISCORD_WEBHOOK_URL_BILLING",
			fmt.Sprintf("💰 User %d upgraded to plan %d", result.UserID, result.NewPlanID),
		)

		if err := CreateUpgradeInvoice.Enqueue(c.Queue(), invoice.CreateUpgradeInput{
			SubscriptionID:  result.SubscriptionID,
			PriceDifference: result.PriceDifference,
			OldPlanID:       result.OldPlanID,
			NewPlanID:       result.NewPlanID,
		}); err != nil {
			return err
		}
	}

	return on_subscription_create.CreateEntitlements.Enqueue(c.Queue(), entitlement.StartUpsertInput{
		UserID: result.UserID,
		PlanID: result.NewPlanID,
	})
}

func upgradeInvoiceHandler(c server.IContext) error {
	var payload invoice.CreateUpgradeInput
	if err := json.Unmarshal(c.Payload(), &payload); err != nil {
		return err
	}

	svc := invoice.NewService(c.DB())
	result, err := svc.CreateUpgrade(&payload)
	if err != nil || !result.ShouldPay {
		return err
	}
	return on_payment.PayInvoice.Enqueue(c.Queue(), map[string]any{
		"invoice_id": result.Invoice.ID,
	})
}
