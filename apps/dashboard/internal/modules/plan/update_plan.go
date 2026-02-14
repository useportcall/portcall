package plan

import (
	"strings"

	quotemodule "github.com/useportcall/portcall/apps/dashboard/internal/modules/quote"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdatePlanRequest struct {
	Name            string `json:"name"`
	Currency        string `json:"currency"`
	TrialPeriodDays int    `json:"trial_period_days"`
	Interval        string `json:"interval"`
	IntervalCount   *int   `json:"interval_count"`
	PlanGroupID     string `json:"plan_group_id"`
	DiscountPct     *int   `json:"discount_pct"`
}

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
	locked, err := quotemodule.HasLockedQuoteForPlan(c, plan.ID)
	if err != nil {
		c.ServerError("Failed to validate quote state", err)
		return
	}
	if locked {
		c.BadRequest("Plan cannot be edited after quote is issued")
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

	if body.DiscountPct != nil {
		plan.DiscountPct = *body.DiscountPct
		plan.DiscountQty = 1
	}

	if err := c.DB().Save(plan); err != nil {
		c.ServerError("Failed to update plan", err)
		return
	}

	planItems := []models.PlanItem{}
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list plan items", err)
		return
	}

	response := new(apix.Plan)
	response.Items = make([]apix.PlanItem, len(planItems))
	for i, pi := range planItems {
		response.Items[i].Set(&pi)
	}

	c.OK(response.Set(plan))
}
