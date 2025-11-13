package plan

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_feature"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_group"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_item"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListPlans(c *routerx.Context) {

	var q []any
	groupID := c.Param("group_id")
	if groupID != "" {
		var group models.PlanGroup
		if err := c.DB().GetForPublicID(c.AppID(), groupID, &group); err != nil {
			c.NotFound("Plan group not found")
			return
		}
		q = []any{"app_id = ? AND plan_group_id = ?", c.AppID(), group.ID}
	} else {
		q = []any{"app_id = ?", c.AppID()}
	}

	plans := []models.Plan{}
	if err := c.DB().ListWithOrder(&plans, "name DESC, created_at DESC", q...); err != nil {
		c.ServerError("Failed to list plans", err)
		return
	}

	response := make([]Plan, len(plans))
	for i, plan := range plans {
		response[i].Set(&plan)
	}

	for i, plan := range plans {
		planItems := []models.PlanItem{}
		if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
			c.ServerError("Failed to list plan items", err)
			return
		}

		if plan.PlanGroupID != nil {
			var planGroup models.PlanGroup
			if err := c.DB().FindForID(*plan.PlanGroupID, &planGroup); err == nil {
				response[i].PlanGroup = (&plan_group.PlanGroup{}).Set(&planGroup)
			}
		}

		response[i].Items = make([]plan_item.PlanItem, len(planItems))
		for j, item := range planItems {
			response[i].Items[j] = plan_item.PlanItem{}
			response[i].Items[j].Set(&item)
		}

		features := []models.PlanFeature{}
		if err := c.DB().List(&features, "plan_id = ?", plan.ID); err != nil {
			c.ServerError("Failed to list plan features", err)
			return
		}

		response[i].Features = make([]any, len(features))
		for j, feature := range features {
			response[i].Features[j] = (&plan_feature.PlanFeature{}).Set(&feature)
		}
	}

	c.OK(response)
}
