package plan

import (
	plan_item "github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_item"
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreatePlan(c *routerx.Context) {
	plan := &models.Plan{
		PublicID:         utils.GenPublicID("plan"),
		AppID:            c.AppID(),
		Name:             "New Plan",
		Status:           "init",
		TrialPeriodDays:  0,
		Interval:         "month",
		IntervalCount:    1,
		Currency:         "USD",
		InvoiceDueByDays: 10}
	if err := c.DB().Create(plan); err != nil {
		c.ServerError("Failed to create plan", err)
		return
	}

	planItem := &models.PlanItem{
		AppID:             c.AppID(),
		PublicID:          utils.GenPublicID("pi"),
		PlanID:            plan.ID,
		PublicTitle:       "New Plan Item",
		PublicDescription: "",
		Quantity:          1,
		UnitAmount:        1000,
		PricingModel:      "fixed"}
	if err := c.DB().Create(planItem); err != nil {
		c.ServerError("Failed to create plan item", err)
		return
	}

	plan.Status = "draft"
	if err := c.DB().Save(plan); err != nil {
		c.ServerError("Failed to update plan status", err)
		return
	}

	response := new(Plan).Set(plan)
	response.Items = append(response.Items, *new(plan_item.PlanItem).Set(planItem))

	c.OK(response)
}
