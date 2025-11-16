package plan

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func PublishPlan(c *routerx.Context) {
	id := c.Param("id")
	response := new(apix.Plan)

	plan := models.Plan{}
	if err := c.DB().GetForPublicID(c.AppID(), id, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	if plan.Name == "" {
		c.BadRequest("Plan does not have a name set")
		return
	}
	if plan.Interval == "" {
		c.BadRequest("Plan does not have a billing cycle interval set")
		return
	}
	if plan.IntervalCount < 1 {
		c.BadRequest("Plan does not have a valid billing cycle count set")
		return
	}

	planItems := []models.PlanItem{}
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list plan items", err)
		return
	}

	for _, pi := range planItems {
		planItem := apix.PlanItem{}
		planItem.Set(&pi)
		if planItem.UnitAmount > 0 || planItem.Tiers != nil {
			plan.IsFree = false
		}
		response.Items = append(response.Items, planItem)
	}

	// Update the plan status to "published"
	plan.Status = "published"
	if err := c.DB().Save(&plan); err != nil {
		c.ServerError("Failed to update plan status", err)
		return
	}

	c.OK(response.Set(&plan))
}
