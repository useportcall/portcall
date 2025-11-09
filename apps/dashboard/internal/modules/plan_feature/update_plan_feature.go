package plan_feature

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func UpdatePlanFeature(c *routerx.Context) {
	id := c.Param("id")

	body := new(UpdatePlanFeatureRequest)
	if err := c.BindJSON(body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	planFeature := &models.PlanFeature{}
	if err := c.DB().GetForPublicID(c.AppID(), id, planFeature); err != nil {
		c.NotFound("Plan feature not found")
		return
	}

	if body.PlanItemID != "" {
		var planItem models.PlanItem
		if err := c.DB().GetForPublicID(c.AppID(), body.PlanItemID, &planItem); err != nil {
			c.NotFound("Plan item not found")
			return
		}
		planFeature.PlanItemID = planItem.ID
	}

	if body.Interval != "" {
		planFeature.Interval = body.Interval
	}

	if body.Quota != 0 {
		planFeature.Quota = body.Quota
	}

	if body.Rollover != nil {
		planFeature.Rollover = *body.Rollover
	}

	if body.FeatureID != nil {
		var feature models.Feature
		if err := c.DB().GetForPublicID(c.AppID(), *body.FeatureID, &feature); err != nil {
			c.NotFound("Feature not found")
			return
		}

		planFeature.FeatureID = feature.ID
	}

	if err := c.DB().Save(planFeature); err != nil {
		c.ServerError("Failed to update plan feature")
		return
	}

	c.OK(new(PlanFeature).Set(planFeature))
}
