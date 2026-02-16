package quotes

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type QuoteDetail struct {
	ID             uint       `json:"id"`
	PublicID       string     `json:"public_id"`
	PublicTitle    string     `json:"public_title"`
	PublicName     string     `json:"public_name"`
	CompanyName    string     `json:"company_name"`
	UserID         *uint      `json:"user_id"`
	UserName       string     `json:"user_name"`
	RecipientEmail string     `json:"recipient_email"`
	PlanID         uint       `json:"plan_id"`
	PlanName       string     `json:"plan_name"`
	Status         string     `json:"status"`
	DaysValid      int        `json:"days_valid"`
	Toc            string     `json:"toc"`
	URL            *string    `json:"url"`
	DirectCheckout bool       `json:"direct_checkout"`
	IssuedAt       *time.Time `json:"issued_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func GetQuote(c *routerx.Context) {
	quoteIDStr := c.Param("quote_id")
	quoteID, err := strconv.ParseUint(quoteIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid quote ID")
		return
	}

	var quote models.Quote
	if err := c.DB().FindForID(uint(quoteID), &quote); err != nil {
		c.NotFound("Quote not found")
		return
	}

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

	result := QuoteDetail{
		ID:             quote.ID,
		PublicID:       quote.PublicID,
		PublicTitle:    quote.PublicTitle,
		PublicName:     quote.PublicName,
		CompanyName:    quote.CompanyName,
		UserID:         quote.UserID,
		UserName:       userName,
		RecipientEmail: quote.RecipientEmail,
		PlanID:         quote.PlanID,
		PlanName:       planName,
		Status:         quote.Status,
		DaysValid:      quote.DaysValid,
		Toc:            quote.Toc,
		URL:            quote.URL,
		DirectCheckout: quote.DirectCheckout,
		IssuedAt:       quote.IssuedAt,
		ExpiresAt:      quote.ExpiresAt,
		AcceptedAt:     quote.AcceptedAt,
		CreatedAt:      quote.CreatedAt,
		UpdatedAt:      quote.UpdatedAt,
	}

	c.OK(result)
}
