package plan_item

import (
	quotemodule "github.com/useportcall/portcall/apps/dashboard/internal/modules/quote"
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreatePlanItemRequest struct {
	PlanID            string `json:"plan_id"`
	PricingModel      string `json:"pricing_model"`
	UnitAmount        int64  `json:"unit_amount"`
	PublicTitle       string `json:"public_title"`
	PublicDescription string `json:"public_description"`
	Interval          string `json:"interval"`       // billing interval: inherit (from plan), week, month, year
	IntervalCount     int    `json:"interval_count"` // number of intervals for the billing cycle
	Quota             int64  `json:"quota"`
	Rollover          int    `json:"rollover"`
}

func CreatePlanItem(c *routerx.Context) {
	body := new(CreatePlanItemRequest)
	if err := c.BindJSON(body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if body.PlanID == "" {
		c.BadRequest("Invalid plan_id")
		return
	}

	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), body.PlanID, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}
	locked, err := quotemodule.HasLockedQuoteForPlan(c, plan.ID)
	if err != nil {
		c.ServerError("Failed to validate quote state", err)
		return
	}
	if locked {
		c.BadRequest("Plan item cannot be added after quote is issued")
		return
	}

	var feature models.Feature
	if err := c.DB().FindFirst(&feature, "app_id = ? AND is_metered = ?", c.AppID(), true); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			c.ServerError("Failed to fetch feature", err)
			return
		}

		feature = models.Feature{
			PublicID:  "tokens",
			AppID:     c.AppID(),
			IsMetered: true,
		}
		if err := c.DB().Create(&feature); err != nil {
			c.ServerError("Failed to create feature", err)
			return
		}
	}

	planItem := &models.PlanItem{
		PublicID:          utils.GenPublicID("pi"),
		PlanID:            plan.ID,
		AppID:             plan.AppID,
		Quantity:          1, // Default quantity, can be changed later
		PricingModel:      body.PricingModel,
		UnitAmount:        body.UnitAmount,
		Tiers:             new([]models.Tier),
		Minimum:           nil,
		Maximum:           nil,
		PublicTitle:       body.PublicTitle,
		PublicDescription: body.PublicDescription,
		Interval:          getItemInterval(body.Interval),
		IntervalCount:     getIntervalCount(body.IntervalCount),
	}
	if err := c.DB().Create(planItem); err != nil {
		c.ServerError("Failed to create plan item", err)
		return
	}

	// For plan features, use inherit if no specific interval provided
	featureInterval := plan.Interval
	if body.Interval != "" {
		featureInterval = body.Interval
	}

	planFeature := models.PlanFeature{
		PublicID:   utils.GenPublicID("pf"),
		PlanID:     plan.ID,
		AppID:      plan.AppID,
		FeatureID:  feature.ID,
		PlanItemID: planItem.ID,
		Interval:   featureInterval,
		Quota:      body.Quota,
		Rollover:   body.Rollover,
	}
	if err := c.DB().Create(&planFeature); err != nil {
		c.ServerError("Failed to create plan feature", err)
		return
	}

	c.OK(new(apix.PlanItem).Set(planItem))
}

// getItemInterval returns the interval or default to "inherit"
func getItemInterval(interval string) string {
	if interval == "" {
		return "inherit"
	}
	return interval
}

// getIntervalCount returns the interval count or default to 1
func getIntervalCount(count int) int {
	if count <= 0 {
		return 1
	}
	return count
}
