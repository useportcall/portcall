package plan

import (
	"strings"

	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_item"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func UpdatePlan(c *routerx.Context) {
	id := c.Param("id")

	var body UpdatePlanRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	plan := &models.Plan{}
	if err := c.DB().GetForPublicID(c.AppID(), id, plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	if body.Name != "" {
		plan.Name = body.Name
	}

	plan.TrialPeriodDays = body.TrialPeriodDays

	if body.Currency != "" {
		plan.Currency = strings.ToUpper(body.Currency)
	}

	if body.Interval != "" {
		plan.Interval = body.Interval
	}

	if body.IntervalCount != nil && *body.IntervalCount > 1 {
		plan.IntervalCount = *body.IntervalCount
	}

	if body.PlanGroupID != "" {
		var planGroup models.PlanGroup
		if err := c.DB().GetForPublicID(c.AppID(), body.PlanGroupID, &planGroup); err != nil {
			c.NotFound("Plan group not found")
			return
		}
		plan.PlanGroupID = &planGroup.ID
	}

	if err := c.DB().Save(plan); err != nil {
		c.ServerError("Failed to update plan")
		return
	}

	planItems := []models.PlanItem{}
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list plan items")
		return
	}

	response := new(Plan)
	response.Items = make([]plan_item.PlanItem, len(planItems))
	for i, pi := range planItems {
		response.Items[i].Set(&pi)
	}

	c.OK(response.Set(plan))
}
