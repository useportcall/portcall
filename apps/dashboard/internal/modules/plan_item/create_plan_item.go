package plan_item

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

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

	var feature models.Feature
	if err := c.DB().FindFirst(&feature, "app_id = ?", c.AppID()); err != nil {
		c.NotFound("Feature not found")
		return
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
	}
	if err := c.DB().Create(planItem); err != nil {
		c.ServerError("Failed to create plan item")
		return
	}

	interval := plan.Interval
	if i := body.Interval; interval != "" {
		interval = i
	}

	planFeature := models.PlanFeature{
		PublicID:   utils.GenPublicID("pf"),
		PlanID:     plan.ID,
		AppID:      plan.AppID,
		FeatureID:  feature.ID,
		PlanItemID: planItem.ID,
		Interval:   interval,
		Quota:      body.Quota,
		Rollover:   body.Rollover,
	}
	if err := c.DB().Create(&planFeature); err != nil {
		c.ServerError("Failed to create plan feature")
		return
	}

	c.OK(new(PlanItem).Set(planItem))
}
