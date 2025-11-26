package quote

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateQuote(c *routerx.Context) {
	plan := &models.Plan{
		PublicID:         utils.GenPublicID("plan"),
		AppID:            c.AppID(),
		Name:             "New Quote",
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

	quote := &models.Quote{
		PublicID:    utils.GenPublicID("quote"),
		AppID:       c.AppID(),
		Status:      "init",
		PublicTitle: "",
		PublicName:  ""}
	if err := c.DB().Create(quote); err != nil {
		c.ServerError("Failed to create quote", err)
		return
	}

	planItem := &models.PlanItem{
		AppID:             c.AppID(),
		PublicID:          utils.GenPublicID("pi"),
		PlanID:            plan.ID,
		PublicTitle:       "",
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

	quote.Status = "draft"
	if err := c.DB().Save(plan); err != nil {
		c.ServerError("Failed to update plan status", err)
		return
	}

	response := new(apix.Plan).Set(plan)
	response.Items = append(response.Items, *new(apix.PlanItem).Set(planItem))

	c.OK(response)
}
