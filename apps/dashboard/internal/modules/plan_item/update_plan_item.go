package plan_item

import (
	quotemodule "github.com/useportcall/portcall/apps/dashboard/internal/modules/quote"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdatePlanItemRequest struct {
	PricingModel      string        `json:"pricing_model"`
	Quantity          *int32        `json:"quantity"`
	UnitAmount        *int64        `json:"unit_amount"`
	Tiers             []models.Tier `json:"tiers"`
	Minimum           *int          `json:"minimum"`
	Maximum           *int          `json:"maximum"`
	PublicTitle       string        `json:"public_title"`
	PublicDescription string        `json:"public_description"`
	Interval          string        `json:"interval"`       // billing interval: inherit, week, month, year
	IntervalCount     *int          `json:"interval_count"` // number of intervals for the billing cycle
}

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
	locked, err := quotemodule.HasLockedQuoteForPlan(c, planItem.PlanID)
	if err != nil {
		c.ServerError("Failed to validate quote state", err)
		return
	}
	if locked {
		c.BadRequest("Plan item cannot be edited after quote is issued")
		return
	}

	if body.PricingModel != "" {
		if body.PricingModel == "fixed" {
			c.BadRequest("cannot change pricing model for fixed plan item")
			return
		}

		if body.PricingModel != "unit" && (planItem.Tiers == nil || len(*planItem.Tiers) == 0) {
			tiers := []models.Tier{}
			tiers = append(tiers, models.Tier{
				Start:  0,
				End:    -1,
				Amount: 1000,
			})
			planItem.Tiers = &tiers
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

	// Update billing interval for the plan item
	if body.Interval != "" {
		// Validate interval value
		switch body.Interval {
		case "inherit", "week", "month", "year":
			planItem.Interval = body.Interval
		default:
			c.BadRequest("invalid interval: must be inherit, week, month, or year")
			return
		}
	}

	if body.IntervalCount != nil && *body.IntervalCount > 0 {
		planItem.IntervalCount = *body.IntervalCount
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
