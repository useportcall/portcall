package plan_feature

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreatePlanFeature(c *routerx.Context) {
	var body CreatePlanFeatureRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if body.PlanID == "" {
		c.BadRequest("'plan_id' property is required")
		return
	}

	if body.FeatureID == "" {
		c.BadRequest("'feature_id' property is required")
		return
	}

	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), body.PlanID, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	var planItem models.PlanItem
	if err := c.DB().FindFirst(&planItem, "plan_id = ? AND pricing_model = 'fixed'", plan.ID); err != nil {
		c.NotFound("Plan item not found")
		return
	}

	var feature models.Feature
	if err := c.DB().GetForPublicID(c.AppID(), body.FeatureID, &feature); err != nil {
		c.NotFound("Feature not found")
		return
	}

	planFeature := &models.PlanFeature{
		PublicID:   dbx.GenPublicID("pf"),
		AppID:      c.AppID(),
		PlanID:     plan.ID,
		PlanItemID: planItem.ID,
		FeatureID:  feature.ID,
		Interval:   plan.Interval,
		Quota:      -1,
	}
	if err := c.DB().Create(planFeature); err != nil {
		c.ServerError("Failed to create plan feature")
		return
	}

	c.OK(new(PlanFeature).Set(planFeature))
}
