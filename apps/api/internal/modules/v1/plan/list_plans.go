package plan

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListPlans(c *routerx.Context) {
	plans := []models.Plan{}
	if err := c.DB().ListWithOrder(&plans, "created_at DESC", "app_id = ?", c.AppID()); err != nil {
		c.ServerError("Failed to list plans", err)
		return
	}

	response := make([]apix.Plan, len(plans))
	for i, p := range plans {
		response[i].Set(&p)

		// Include plan items
		planItems := []models.PlanItem{}
		if err := c.DB().List(&planItems, "plan_id = ?", p.ID); err != nil {
			c.ServerError("Failed to list plan items", err)
			return
		}

		response[i].Items = make([]apix.PlanItem, len(planItems))
		for j, item := range planItems {
			response[i].Items[j] = apix.PlanItem{}
			response[i].Items[j].Set(&item)
		}

		// Include plan features
		planFeatures := []models.PlanFeature{}
		if err := c.DB().List(&planFeatures, "plan_id = ?", p.ID); err != nil {
			c.ServerError("Failed to list plan features", err)
			return
		}

		response[i].Features = make([]any, len(planFeatures))
		for j, pf := range planFeatures {
			response[i].Features[j] = (&apix.PlanFeature{}).Set(&pf)
		}
	}

	c.OK(response)
}
