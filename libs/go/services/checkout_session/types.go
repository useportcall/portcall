package checkout_session

import "github.com/useportcall/portcall/libs/go/dbx/models"

// CreateInput is the input for creating a checkout session.
type CreateInput struct {
	AppID                 uint
	PlanID                string // public ID
	UserID                string // public ID
	CancelURL             string
	RedirectURL           string
	RequireBillingAddress bool
}

// CreateResult is the result of creating a checkout session.
type CreateResult struct {
	Session     *models.CheckoutSession
	CheckoutURL string
}

// ResolvePayload is the payload for resolving a checkout session.
type ResolvePayload struct {
	ExternalSessionID       string `json:"external_session_id"`
	ExternalPaymentMethodID string `json:"external_payment_method_id"`
}

// ResolveResult is the result of resolving a checkout session.
type ResolveResult struct {
	Session                 *models.CheckoutSession
	ExternalPaymentMethodID string
	Skipped                 bool
}
