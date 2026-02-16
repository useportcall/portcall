package subscriptions

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type SubscriptionListItem struct {
	ID              uint      `json:"id"`
	PublicID        string    `json:"public_id"`
	UserID          uint      `json:"user_id"`
	UserName        string    `json:"user_name"`
	UserEmail       string    `json:"user_email"`
	PlanID          *uint     `json:"plan_id"`
	PlanName        string    `json:"plan_name"`
	Status          string    `json:"status"`
	BillingInterval string    `json:"billing_interval"`
	Currency        string    `json:"currency"`
	NextResetAt     time.Time `json:"next_reset_at"`
	ItemCount       int64     `json:"item_count"`
}

func ListSubscriptions(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid app ID")
		return
	}

	var subscriptions []models.Subscription
	if err := c.DB().List(&subscriptions, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list subscriptions", err)
		return
	}

	result := make([]SubscriptionListItem, len(subscriptions))
	for i, sub := range subscriptions {
		// Get user info
		var user models.User
		userName := ""
		userEmail := ""
		if err := c.DB().FindForID(sub.UserID, &user); err == nil {
			userName = user.Name
			userEmail = user.Email
		}

		// Get plan name
		planName := ""
		if sub.PlanID != nil {
			var plan models.Plan
			if err := c.DB().FindForID(*sub.PlanID, &plan); err == nil {
				planName = plan.Name
			}
		}

		var itemCount int64
		c.DB().Count(&itemCount, models.SubscriptionItem{}, "subscription_id = ?", sub.ID)

		result[i] = SubscriptionListItem{
			ID:              sub.ID,
			PublicID:        sub.PublicID,
			UserID:          sub.UserID,
			UserName:        userName,
			UserEmail:       userEmail,
			PlanID:          sub.PlanID,
			PlanName:        planName,
			Status:          sub.Status,
			BillingInterval: sub.BillingInterval,
			Currency:        sub.Currency,
			NextResetAt:     sub.NextResetAt,
			ItemCount:       itemCount,
		}
	}

	c.OK(result)
}
