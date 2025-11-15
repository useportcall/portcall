package plan_item

import (
	"strconv"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListPlanItems(c *routerx.Context) {
	planID := c.Query("plan_id")
	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), planID, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	query := "plan_id = ?"
	args := []any{plan.ID}

	if pricingModel := c.Query("pricing_model"); pricingModel != "" {
		query += " AND pricing_model = ?"
		args = append(args, pricingModel)
	} else if isMetered := c.Query("is_metered"); isMetered != "" {
		query += " AND pricing_model <> 'fixed'"
	}

	planItems := []models.PlanItem{}
	if err := c.DB().ListWithOrderAndLimit(&planItems, "created_at DESC", limit, append([]any{query}, args...)...); err != nil {
		c.ServerError("Failed to list plan items", err)
		return
	}

	response := make([]apix.PlanItem, len(planItems))
	for i, pi := range planItems {
		response[i] = *new(apix.PlanItem)
		response[i].Set(&pi)

		planFeatures := []models.PlanFeature{}
		if err := c.DB().List(&planFeatures, "plan_item_id = ?", pi.ID); err != nil {
			c.ServerError("Failed to list plan features", err)
			return
		}

		for _, pf := range planFeatures {
			planFeature := apix.PlanFeature{}
			planFeature.Set(&pf)

			var feature models.Feature
			if err := c.DB().FindForID(pf.FeatureID, &feature); err != nil {
				c.ServerError("Failed to find feature", err)
				return
			}
			planFeature.Feature = map[string]any{"id": feature.PublicID}
			response[i].Features = append(response[i].Features, planFeature)
		}
	}

	c.OK(response)
}
