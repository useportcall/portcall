package handlers

import (
	"encoding/json"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type ProcessPlanSwitchPayload struct {
	OldPlanID      uint `json:"old_plan_id"`
	NewPlanID      uint `json:"new_plan_id"`
	SubscriptionID uint `json:"subscription_id"`
}

func ProcessPlanSwitch(c server.IContext) error {
	var p ProcessPlanSwitchPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var oldFixedPlanItem models.PlanItem
	if err := c.DB().FindFirst(&oldFixedPlanItem, "plan_id = ? AND pricing_model = ?", p.OldPlanID, "fixed"); err != nil {
		return err
	}

	var newFixedPlanItem models.PlanItem
	if err := c.DB().FindFirst(&newFixedPlanItem, "plan_id = ? AND pricing_model = ?", p.NewPlanID, "fixed"); err != nil {
		return err
	}

	if newFixedPlanItem.UnitAmount > oldFixedPlanItem.UnitAmount {
		if err := c.Queue().Enqueue(
			"create_invoice",
			map[string]any{"subscription_id": p.SubscriptionID},
			"billing_queue",
		); err != nil {
			return err
		}
	}

	return nil
}
