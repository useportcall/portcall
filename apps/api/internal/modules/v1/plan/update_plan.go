package plan

import (
	"strings"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdatePlanRequest struct {
	Name            string `json:"name"`
	Currency        string `json:"currency"`
	TrialPeriodDays *int   `json:"trial_period_days"`
	Interval        string `json:"interval"`
	IntervalCount   *int   `json:"interval_count"`
	UnitAmount      *int64 `json:"unit_amount"`
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

	if body.Name != "" {
		plan.Name = body.Name
	}

	if body.TrialPeriodDays != nil {
		plan.TrialPeriodDays = *body.TrialPeriodDays
	}

	if body.Currency != "" {
		plan.Currency = strings.ToUpper(body.Currency)
	}

	if body.Interval != "" {
		plan.Interval = body.Interval
	}

	if body.IntervalCount != nil && *body.IntervalCount >= 1 {
		plan.IntervalCount = *body.IntervalCount
	}

	if err := c.DB().Save(plan); err != nil {
		c.ServerError("Failed to update plan", err)
		return
	}

	// Update the fixed plan item's unit amount if provided
	if body.UnitAmount != nil {
		var planItem models.PlanItem
		if err := c.DB().FindFirst(&planItem, "plan_id = ? AND pricing_model = 'fixed'", plan.ID); err == nil {
			planItem.UnitAmount = *body.UnitAmount
			if err := c.DB().Save(&planItem); err != nil {
				c.ServerError("Failed to update plan item", err)
				return
			}
		}
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

	features, ok := loadPlanFeatures(c, plan.ID)
	if !ok {
		return
	}
	response.Features = features

	c.OK(response.Set(plan))
}
