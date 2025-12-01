package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Quote struct {
	ID                    string     `json:"id"`
	Status                string     `json:"status"`
	Email                 string     `json:"email"` // email of the user to whom the quote is sent
	CompanyName           string     `json:"company_name"`
	DaysValid             int        `json:"days_valid"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	IssuedAt              *time.Time `json:"issued_at"`
	ExpiresAt             *time.Time `json:"expires_at"`
	AcceptedAt            *time.Time `json:"accepted_at"`
	VoidedAt              *time.Time `json:"voided_at"`
	Plan                  any        `json:"plan"`
	User                  any        `json:"user"`
	DirectCheckoutEnabled bool       `json:"direct_checkout_enabled"`
}

func (q *Quote) Set(quote *models.Quote) *Quote {
	q.ID = quote.PublicID
	q.Status = quote.Status
	q.CompanyName = quote.CompanyName
	q.DaysValid = quote.DaysValid
	q.CreatedAt = quote.CreatedAt
	q.UpdatedAt = quote.UpdatedAt
	q.ExpiresAt = quote.ExpiresAt
	q.AcceptedAt = quote.AcceptedAt
	q.VoidedAt = quote.VoidedAt
	q.DirectCheckoutEnabled = quote.DirectCheckout

	if quote.IssuedAt != nil {
		q.IssuedAt = quote.IssuedAt
	}

	if quote.Plan.ID != 0 {
		q.Plan = new(Plan).Set(&quote.Plan)
	}

	if quote.User != nil {
		q.User = new(User).Set(quote.User)
	}

	return q
}
