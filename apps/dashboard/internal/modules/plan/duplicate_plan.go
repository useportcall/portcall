package plan

import (
	plan_item "github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_item"
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DuplicatePlan(c *routerx.Context) {
	id := c.Param("id")

	plan := &models.Plan{}
	if err := c.DB().GetForPublicID(c.AppID(), id, plan); err != nil {
		c.NotFound("Original plan not found")
		return
	}

	var originalPlanItems []models.PlanItem
	if err := c.DB().List(&originalPlanItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Failed to list original plan items", err)
		return
	}

	// Copy relevant fields from the original plan
	plan.PublicID = utils.GenPublicID("plan")
	plan.ID = 0 // Reset ID for new plan
	plan.Status = "init"
	if err := c.DB().Create(plan); err != nil {
		c.ServerError("Failed to create plan", err)
		return
	}

	result := new(Plan)
	result.Items = make([]plan_item.PlanItem, len(originalPlanItems))
	for i, item := range originalPlanItems {
		// plan features
		var planFeatures []models.PlanFeature
		if err := c.DB().List(&planFeatures, "plan_item_id = ?", item.ID); err != nil {
			c.ServerError("Failed to list plan features", err)
			return
		}

		item.PublicID = utils.GenPublicID("plan_item")
		item.ID = 0           // Reset ID for new item
		item.PlanID = plan.ID // Associate with new plan
		if err := c.DB().Create(&item); err != nil {
			c.ServerError("Failed to create plan item", err)
			return
		}

		for _, pf := range planFeatures {
			pf.PublicID = utils.GenPublicID("plan_feature")
			pf.ID = 0               // Reset ID for new feature
			pf.PlanID = plan.ID     // Associate with new plan
			pf.PlanItemID = item.ID // Associate with new plan item
			if err := c.DB().Create(&pf); err != nil {
				c.ServerError("Failed to create plan feature", err)
				return
			}
		}

		result.Items[i] = *(&plan_item.PlanItem{}).Set(&item)
	}

	plan.Status = "draft"
	if err := c.DB().Save(plan); err != nil {
		c.ServerError("Failed to update plan status", err)
		return
	}

	c.OK(result.Set(plan))
}
