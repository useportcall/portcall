package plans

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type PlanDetail struct {
	ID                uint         `json:"id"`
	PublicID          string       `json:"public_id"`
	Name              string       `json:"name"`
	Status            string       `json:"status"`
	Interval          string       `json:"interval"`
	IntervalCount     int          `json:"interval_count"`
	Currency          string       `json:"currency"`
	IsFree            bool         `json:"is_free"`
	TrialPeriodDays   int          `json:"trial_period_days"`
	SubscriptionCount int64        `json:"subscription_count"`
	Items             []ItemDetail `json:"items"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type ItemDetail struct {
	ID           uint   `json:"id"`
	PublicID     string `json:"public_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Quantity     int32  `json:"quantity"`
	UnitAmount   int64  `json:"unit_amount"`
	PricingModel string `json:"pricing_model"`
}

func GetPlan(c *routerx.Context) {
	planIDStr := c.Param("plan_id")
	planID, err := strconv.ParseUint(planIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid plan ID")
		return
	}

	var plan models.Plan
	if err := c.DB().FindForID(uint(planID), &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	// Get items
	var items []models.PlanItem
	c.DB().List(&items, "plan_id = ?", plan.ID)

	itemDetails := make([]ItemDetail, len(items))
	for i, item := range items {
		itemDetails[i] = ItemDetail{
			ID:           item.ID,
			PublicID:     item.PublicID,
			Title:        item.PublicTitle,
			Description:  item.PublicDescription,
			Quantity:     item.Quantity,
			UnitAmount:   item.UnitAmount,
			PricingModel: item.PricingModel,
		}
	}

	var subCount int64
	c.DB().Count(&subCount, models.Subscription{}, "plan_id = ?", plan.ID)

	result := PlanDetail{
		ID:                plan.ID,
		PublicID:          plan.PublicID,
		Name:              plan.Name,
		Status:            plan.Status,
		Interval:          plan.Interval,
		IntervalCount:     plan.IntervalCount,
		Currency:          plan.Currency,
		IsFree:            plan.IsFree,
		TrialPeriodDays:   plan.TrialPeriodDays,
		SubscriptionCount: subCount,
		Items:             itemDetails,
		CreatedAt:         plan.CreatedAt,
		UpdatedAt:         plan.UpdatedAt,
	}

	c.OK(result)
}
