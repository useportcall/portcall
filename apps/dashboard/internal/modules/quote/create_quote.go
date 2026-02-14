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

	address := new(models.Address)
	// ensure PublicID is set to avoid unique index collisions
	address.PublicID = utils.GenPublicID("address")
	address.AppID = c.AppID()
	if err := c.DB().Create(address); err != nil {
		c.ServerError("Failed to create recipient address", err)
		return
	}

	quote := &models.Quote{
		PublicID:           utils.GenPublicID("quote"),
		AppID:              c.AppID(),
		UserID:             nil,
		PlanID:             plan.ID,
		RecipientAddressID: address.ID,
		Status:             "init",
		PublicTitle:        "",
		PublicName:         ""}
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

	quote.Plan = *plan

	quote.Status = "draft"
	if err := c.DB().Save(quote); err != nil {
		c.ServerError("Failed to update quote status", err)
		return
	}

	response := new(apix.Quote).Set(quote)

	c.OK(response)
}
