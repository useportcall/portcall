package plans

import (
	"strconv"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type PlanListItem struct {
	ID                uint   `json:"id"`
	PublicID          string `json:"public_id"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	Interval          string `json:"interval"`
	IntervalCount     int    `json:"interval_count"`
	Currency          string `json:"currency"`
	IsFree            bool   `json:"is_free"`
	SubscriptionCount int64  `json:"subscription_count"`
	ItemCount         int64  `json:"item_count"`
}

func ListPlans(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid app ID")
		return
	}

	var plans []models.Plan
	if err := c.DB().List(&plans, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list plans", err)
		return
	}

	result := make([]PlanListItem, len(plans))
	for i, plan := range plans {
		var subCount, itemCount int64
		c.DB().Count(&subCount, models.Subscription{}, "plan_id = ?", plan.ID)
		c.DB().Count(&itemCount, models.PlanItem{}, "plan_id = ?", plan.ID)

		result[i] = PlanListItem{
			ID:                plan.ID,
			PublicID:          plan.PublicID,
			Name:              plan.Name,
			Status:            plan.Status,
			Interval:          plan.Interval,
			IntervalCount:     plan.IntervalCount,
			Currency:          plan.Currency,
			IsFree:            plan.IsFree,
			SubscriptionCount: subCount,
			ItemCount:         itemCount,
		}
	}

	c.OK(result)
}
