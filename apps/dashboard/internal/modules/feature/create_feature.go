package feature

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateFeature(c *routerx.Context) {
	var body CreateFeatureRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	feature := models.Feature{}

	if body.FeatureID != "" {
		feature.PublicID = body.FeatureID
	} else {
		c.BadRequest("feature_id is required")
		return
	}

	feature.IsMetered = body.IsMetered
	feature.AppID = c.AppID()

	if err := c.DB().Create(&feature); err != nil {
		c.ServerError("Failed to create feature")
		return
	}

	if body.PlanID != "" {
		plan := &models.Plan{}
		if err := c.DB().GetForPublicID(c.AppID(), body.PlanID, plan); err != nil {
			c.ServerError("Failed to get plan")
			return
		}

		var planItem models.PlanItem
		if err := c.DB().FindFirst(&planItem, "plan_id = ? AND pricing_model = 'fixed'", plan.ID); err != nil {
			c.ServerError("Failed to get plan item")
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
	}

	if body.PlanFeatureID != "" {
		var planFeature models.PlanFeature
		if err := c.DB().GetForPublicID(c.AppID(), body.PlanFeatureID, &planFeature); err != nil {
			c.ServerError("Failed to get plan feature")
			return
		}

		planFeature.FeatureID = feature.ID
		if err := c.DB().Save(planFeature); err != nil {
			c.ServerError("Failed to update plan feature")
			return
		}
	}

	c.OK(new(Feature).Set(&feature))
}
