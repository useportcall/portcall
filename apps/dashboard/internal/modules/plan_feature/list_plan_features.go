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
	pricingModel := c.Query("pricing_model")

	if planID != "" {
		var plan models.Plan
		if err := c.DB().GetForPublicID(c.AppID(), planID, &plan); err != nil {
			c.NotFound("Plan not found")
			return
		}
		query += " AND plan_id = ?"
		values = append(values, plan.ID)
	}

	if planItemID != "" {
		var planItem models.PlanItem
		if err := c.DB().GetForPublicID(c.AppID(), planItemID, &planItem); err != nil {
			c.NotFound("Plan item not found")
			return
		}
		query += " AND plan_item_id = ?"
		values = append(values, planItem.ID)
	}

	if pricingModel != "" {
		query += " AND pricing_model = ?"
		values = append(values, pricingModel)
	}

	pfQuery := append([]any{query}, values...)

	planFeatures := []models.PlanFeature{}
	if err := c.DB().List(&planFeatures, pfQuery...); err != nil {
		c.ServerError("Failed to list plan features")
		return
	}

	response := []PlanFeature{}
	for _, pf := range planFeatures {
		data := PlanFeature{}
		data.Set(&pf)

		var planItem models.PlanItem
		if err := c.DB().FindForID(pf.PlanItemID, &planItem); err == nil {
			data.PlanItem = planItem.PublicID
		}

		var f models.Feature
		if err := c.DB().FindForID(pf.FeatureID, &f); err == nil {
			data.Feature = new(feature.Feature).Set(&f)
		}

		response = append(response, data)
	}

	c.OK(response)
}
