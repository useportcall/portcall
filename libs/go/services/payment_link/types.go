package payment_link

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type CreateInput struct {
	AppID                 uint
	PlanID                string
	UserID                string
	UserEmail             string
	UserName              string
	CancelURL             string
	RedirectURL           string
	ExpiresAt             *time.Time
	RequireBillingAddress bool
}

type CreateResult struct {
	PaymentLink *models.PaymentLink
	URL         string
}

type RedeemInput struct {
	ID    string
	Token string
}

type RedeemResult struct {
	PaymentLink *models.PaymentLink
	Session     *models.CheckoutSession
	CheckoutURL string
}
