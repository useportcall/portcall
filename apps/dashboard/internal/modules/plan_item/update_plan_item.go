package plan_item

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func UpdatePlanItem(c *routerx.Context) {
	body := new(UpdatePlanItemRequest)
	if err := c.BindJSON(body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	id := c.Param("id")
	planItem := &models.PlanItem{}
	if err := c.DB().GetForPublicID(c.AppID(), id, planItem); err != nil {
		c.NotFound("Plan item not found")
		return
	}

	if body.PricingModel != "" {
		if body.PricingModel == "fixed" {
			c.BadRequest("cannot change pricing model for fixed plan item")
			return
		}
		planItem.PricingModel = body.PricingModel
	}

	if body.Quantity != nil {
		planItem.Quantity = *body.Quantity
	}

	if body.UnitAmount != nil {
		planItem.UnitAmount = *body.UnitAmount
	}

	if body.Tiers != nil {
		planItem.Tiers = &body.Tiers
	}

	if body.Minimum != nil {
		planItem.Minimum = body.Minimum
	}

	if body.Maximum != nil {
		planItem.Maximum = body.Maximum
	}

	if body.PublicTitle != "" {
		planItem.PublicTitle = body.PublicTitle
	}

	if body.PublicDescription != "" {
		planItem.PublicDescription = body.PublicDescription
	}

	if (planItem.Tiers == nil || len(*planItem.Tiers) == 0) && planItem.PricingModel == "tiered" {
		c.BadRequest("tiers must be set for tiered pricing model")
		return
	}

	if err := c.DB().Save(planItem); err != nil {
		c.ServerError("Failed to update plan item", err)
		return
	}

	c.OK(new(apix.PlanItem).Set(planItem))
}
