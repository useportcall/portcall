package quotes

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type QuoteListItem struct {
	ID             uint      `json:"id"`
	PublicID       string    `json:"public_id"`
	PublicTitle    string    `json:"public_title"`
	UserID         *uint     `json:"user_id"`
	UserName       string    `json:"user_name"`
	RecipientEmail string    `json:"recipient_email"`
	PlanID         uint      `json:"plan_id"`
	PlanName       string    `json:"plan_name"`
	Status         string    `json:"status"`
	DaysValid      int       `json:"days_valid"`
	CreatedAt      time.Time `json:"created_at"`
}

func ListQuotes(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid app ID")
		return
	}

	var quotes []models.Quote
	if err := c.DB().List(&quotes, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list quotes", err)
		return
	}

	result := make([]QuoteListItem, len(quotes))
	for i, quote := range quotes {
		// Get user name
		userName := ""
		if quote.UserID != nil {
			var user models.User
			if err := c.DB().FindForID(*quote.UserID, &user); err == nil {
				userName = user.Name
			}
		}

		// Get plan name
		planName := ""
		if quote.PlanID != 0 {
			var plan models.Plan
			if err := c.DB().FindForID(quote.PlanID, &plan); err == nil {
				planName = plan.Name
			}
		}

		result[i] = QuoteListItem{
			ID:             quote.ID,
			PublicID:       quote.PublicID,
			PublicTitle:    quote.PublicTitle,
			UserID:         quote.UserID,
			UserName:       userName,
			RecipientEmail: quote.RecipientEmail,
			PlanID:         quote.PlanID,
			PlanName:       planName,
			Status:         quote.Status,
			DaysValid:      quote.DaysValid,
			CreatedAt:      quote.CreatedAt,
		}
	}

	c.OK(result)
}
