package plan

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreatePlanRequest struct {
	Name            string `json:"name" binding:"required"`
	Currency        string `json:"currency"`
	Interval        string `json:"interval"`
	IntervalCount   int    `json:"interval_count"`
	TrialPeriodDays int    `json:"trial_period_days"`
	UnitAmount      int64  `json:"unit_amount"`
}

func CreatePlan(c *routerx.Context) {
	var body CreatePlanRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body: 'name' is required")
		return
	}

	// Set defaults
	currency := body.Currency
	if currency == "" {
		currency = "USD"
	}

	interval := body.Interval
	if interval == "" {
		interval = "month"
	}

	intervalCount := body.IntervalCount
	if intervalCount < 1 {
		intervalCount = 1
	}

	plan := &models.Plan{
		PublicID:         dbx.GenPublicID("plan"),
		AppID:            c.AppID(),
		Name:             body.Name,
		Status:           "init",
		TrialPeriodDays:  body.TrialPeriodDays,
		Interval:         interval,
		IntervalCount:    intervalCount,
		Currency:         currency,
		InvoiceDueByDays: 10,
	}
	if err := c.DB().Create(plan); err != nil {
		c.ServerError("Failed to create plan", err)
		return
	}

	// Create a default fixed plan item
	unitAmount := body.UnitAmount
	if unitAmount < 0 {
		unitAmount = 0
	}

	planItem := &models.PlanItem{
		AppID:             c.AppID(),
		PublicID:          dbx.GenPublicID("pi"),
		PlanID:            plan.ID,
		PublicTitle:       "",
		PublicDescription: "",
		Quantity:          1,
		UnitAmount:        unitAmount,
		PricingModel:      "fixed",
	}
	if err := c.DB().Create(planItem); err != nil {
		c.ServerError("Failed to create plan item", err)
		return
	}

	plan.Status = "draft"
	if err := c.DB().Save(plan); err != nil {
		c.ServerError("Failed to update plan status", err)
		return
	}

	response := new(apix.Plan).Set(plan)
	response.Items = append(response.Items, *new(apix.PlanItem).Set(planItem))

	c.OK(response)
}
