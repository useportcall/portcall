package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type CreateSubscriptionItemsPayload struct {
	SubscriptionID uint `json:"subscription_id"`
	PlanID         uint `json:"plan_id"`
}

func CreateSubscriptionItems(c server.IContext) error {
	var p CreateSubscriptionItemsPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var planItems []models.PlanItem
	if err := c.DB().List(&planItems, "plan_id = ?", p.PlanID); err != nil {
		return err
	}

	for _, pi := range planItems {
		subscriptionItem := models.SubscriptionItem{
			PublicID:       dbx.GenPublicID("si"),
			PlanItemID:     &pi.ID,
			Quantity:       pi.Quantity,
			AppID:          pi.AppID,
			UnitAmount:     pi.UnitAmount,
			PricingModel:   pi.PricingModel,
			Tiers:          pi.Tiers,
			Maximum:        pi.Maximum,
			Minimum:        pi.Minimum,
			Title:          pi.PublicTitle,
			Description:    pi.PublicDescription,
			SubscriptionID: p.SubscriptionID,
		}
		if err := c.DB().Create(&subscriptionItem); err != nil {
			return fmt.Errorf("failed to create subscription item for plan item %s: %v", pi.PublicID, err)
		}
	}

	if err := c.Queue().Enqueue("create_invoice", map[string]any{"subscription_id": p.SubscriptionID}, "billing_queue"); err != nil {
		return err
	}

	return nil
}
