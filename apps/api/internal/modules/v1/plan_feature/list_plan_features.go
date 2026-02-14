package plan_feature

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListPlanFeatures(c *routerx.Context) {
	planID := c.Query("plan_id")
	if planID == "" {
		c.BadRequest("'plan_id' query parameter is required")
		return
	}

	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), planID, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	planFeatures := []models.PlanFeature{}
	if err := c.DB().List(&planFeatures, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list plan features", err)
		return
	}

	response := make([]apix.PlanFeature, len(planFeatures))
	for i, pf := range planFeatures {
		response[i].Set(&pf)
		var feature models.Feature
		if err := c.DB().FindForID(pf.FeatureID, &feature); err == nil {
			response[i].Feature = feature.PublicID
		}
	}

	c.OK(response)
}
