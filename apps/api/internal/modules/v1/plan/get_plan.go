package plan

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetPlan(c *routerx.Context) {
	id := c.Param("id")
	response := new(apix.Plan)

	plan := models.Plan{}
	if err := c.DB().GetForPublicID(c.AppID(), id, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	if plan.PlanGroupID != nil {
		var planGroup models.PlanGroup
		if err := c.DB().FindForID(*plan.PlanGroupID, &planGroup); err == nil {
			response.PlanGroup = new(apix.PlanGroup).Set(&planGroup)
		}
	}

	planItems := []models.PlanItem{}
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list plan items", err)
		return
	}

	for _, item := range planItems {
		planItem := apix.PlanItem{}
		planItem.Set(&item)
		response.Items = append(response.Items, planItem)
	}

	planFeatures := []models.PlanFeature{}
	if err := c.DB().List(&planFeatures, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list plan features", err)
		return
	}

	for _, pf := range planFeatures {
		planFeature := apix.PlanFeature{}
		planFeature.Set(&pf)

		var feature models.Feature
		if err := c.DB().FindForID(pf.FeatureID, &feature); err == nil {
			planFeature.Feature = feature.PublicID
		}

		response.Features = append(response.Features, planFeature)
	}

	c.OK(response.Set(&plan))
}
