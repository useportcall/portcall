package plan

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_group"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_item"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetPlan(c *routerx.Context) {
	id := c.Param("id")
	response := new(Plan)

	plan := models.Plan{}
	if err := c.DB().GetForPublicID(c.AppID(), id, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	if plan.PlanGroupID != nil {
		var planGroup models.PlanGroup
		if err := c.DB().FindForID(*plan.PlanGroupID, &planGroup); err == nil {
			response.PlanGroup = new(plan_group.PlanGroup).Set(&planGroup)
		}
	}

	planItems := []models.PlanItem{}
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list plan items")
		return
	}

	for _, item := range planItems {
		planItem := plan_item.PlanItem{}
		planItem.Set(&item)
		response.Items = append(response.Items, planItem)
	}

	c.OK(response.Set(&plan))
}
