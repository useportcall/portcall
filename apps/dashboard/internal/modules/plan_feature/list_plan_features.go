package plan_feature

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/feature"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListPlanFeatures(c *routerx.Context) {
	query := "app_id = ?"
	values := []any{c.AppID()}

	planID := c.Query("plan_id")
	planItemID := c.Query("plan_item_id")

	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), planID, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}
	query += " AND plan_id = ?"
	values = append(values, plan.ID)

	var planItem models.PlanItem

	if planItemID != "" {
		if err := c.DB().GetForPublicID(c.AppID(), planItemID, &planItem); err != nil {
			c.NotFound("Plan item not found")
			return
		}
	} else {
		if err := c.DB().FindFirst(&planItem, "plan_id = ? AND pricing_model = ?", plan.ID, "fixed"); err != nil {
			c.NotFound("Plan item not found")
			return
		}
	}

	query += " AND plan_item_id = ?"
	values = append(values, planItem.ID)

	conds := []any{query}
	conds = append(conds, values...)

	planFeatures := []models.PlanFeature{}
	if err := c.DB().List(&planFeatures, conds...); err != nil {
		c.ServerError("Failed to list plan features", err)
		return
	}

	response := make([]PlanFeature, len(planFeatures))
	for i, pf := range planFeatures {
		data := new(PlanFeature)
		data.Set(&pf)

		var planItem models.PlanItem
		if err := c.DB().FindForID(pf.PlanItemID, &planItem); err != nil {
			c.ServerError("Failed to find plan item", err)
			return
		}

		data.PlanItem = planItem.PublicID

		var f models.Feature
		if err := c.DB().FindForID(pf.FeatureID, &f); err != nil {
			c.ServerError("Failed to find feature", err)
			return
		}

		data.Feature = new(feature.Feature).Set(&f)

		response[i] = *data
	}

	c.OK(response)
}
