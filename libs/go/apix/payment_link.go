package apix

import (
	"fmt"
	"os"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type PaymentLink struct {
	ID                    string    `json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	ExpiresAt             time.Time `json:"expires_at"`
	Status                string    `json:"status"`
	PlanID                string    `json:"plan_id"`
	UserID                string    `json:"user_id"`
	RedirectURL           *string   `json:"redirect_url"`
	CancelURL             *string   `json:"cancel_url"`
	RequireBillingAddress bool      `json:"require_billing_address"`
	URL                   string    `json:"url"`
}

func (pl *PaymentLink) Set(link *models.PaymentLink) *PaymentLink {
	pl.ID = link.PublicID
	pl.CreatedAt = link.CreatedAt
	pl.ExpiresAt = link.ExpiresAt
	pl.Status = link.Status
	pl.RedirectURL = link.RedirectURL
	pl.CancelURL = link.CancelURL
	pl.RequireBillingAddress = link.RequireBillingAddress
	pl.PlanID = link.Plan.PublicID
	pl.UserID = link.User.PublicID

	checkoutURL := os.Getenv("CHECKOUT_URL")
	if checkoutURL != "" {
		pl.URL = fmt.Sprintf("%s?pl=%s", checkoutURL, link.PublicID)
	}

	return pl
}
