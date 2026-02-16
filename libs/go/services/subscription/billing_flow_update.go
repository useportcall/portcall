package subscription

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/services/invoice"
)

func runUpdateFlow(input *BillingFlowInput, subscriptionID uint) error {
	svc := NewService(input.DB)
	updateResult, err := svc.Update(&UpdateInput{
		SubscriptionID: subscriptionID,
		PlanID:         input.PlanID,
		AppID:          input.AppID,
	})
	if err != nil {
		return err
	}

	itemsResult, err := svc.CreateItems(&CreateItemsInput{
		SubscriptionID: subscriptionID,
		PlanID:         updateResult.NewPlanID,
	})
	if err != nil {
		return err
	}

	switchResult, err := svc.PlanSwitch(&PlanSwitchInput{
		OldPlanID:      updateResult.OldPlanID,
		NewPlanID:      updateResult.NewPlanID,
		SubscriptionID: subscriptionID,
	})
	if err != nil {
		return err
	}
	if switchResult.IsUpgrade && input.NotifyUpgrade != nil {
		input.NotifyUpgrade(switchResult.UserID, switchResult.NewPlanID)
	}

	invoiceID, err := createStandardInvoice(input.DB, subscriptionID, itemsResult.ItemCount)
	if err != nil {
		return err
	}

	upgradeInvoiceID, err := createUpgradeInvoice(input.DB, switchResult)
	if err != nil {
		return err
	}

	if err := syncEntitlements(input.DB, switchResult.UserID, switchResult.NewPlanID); err != nil {
		return err
	}
	if invoiceID != 0 {
		if err := payAndResolveInvoice(input, invoiceID); err != nil {
			return err
		}
	}
	if upgradeInvoiceID == 0 {
		return nil
	}
	return payAndResolveInvoice(input, upgradeInvoiceID)
}

func createUpgradeInvoice(db dbx.IORM, result *PlanSwitchResult) (uint, error) {
	if !result.IsUpgrade {
		return 0, nil
	}
	svc := invoice.NewService(db)
	upgrade, err := svc.CreateUpgrade(&invoice.CreateUpgradeInput{
		SubscriptionID:  result.SubscriptionID,
		PriceDifference: result.PriceDifference,
		OldPlanID:       result.OldPlanID,
		NewPlanID:       result.NewPlanID,
	})
	if err != nil || !upgrade.ShouldPay {
		return 0, err
	}
	return upgrade.Invoice.ID, nil
}
