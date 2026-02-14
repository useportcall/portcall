package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Quote struct {
	ID                    string     `json:"id"`
	Status                string     `json:"status"`
	URL                   *string    `json:"url"`
	SignatureURL          *string    `json:"signature_url"`
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
	RecipientAddress      any        `json:"recipient_address"`
	RecipientEmail        string     `json:"recipient_email"`
	RecipientName         string     `json:"recipient_name"`
	RecipientTitle        string     `json:"recipient_title"`
	User                  any        `json:"user"`
	DirectCheckoutEnabled bool       `json:"direct_checkout_enabled"`
	ToC                   string     `json:"toc"` // terms and conditions
	PreparedByEmail       string     `json:"prepared_by_email"`
}

func (q *Quote) Set(quote *models.Quote) *Quote {
	q.ID = quote.PublicID
	q.Status = quote.Status
	q.URL = quote.URL
	q.SignatureURL = quote.SignatureURL
	q.CompanyName = quote.CompanyName
	q.DaysValid = quote.DaysValid
	q.CreatedAt = quote.CreatedAt
	q.UpdatedAt = quote.UpdatedAt
	q.ExpiresAt = quote.ExpiresAt
	q.AcceptedAt = quote.AcceptedAt
	q.VoidedAt = quote.VoidedAt
	q.DirectCheckoutEnabled = quote.DirectCheckout
	q.ToC = quote.Toc
	q.RecipientEmail = quote.RecipientEmail
	q.RecipientName = quote.PublicName
	q.RecipientTitle = quote.PublicTitle
	q.PreparedByEmail = quote.PreparedByEmail

	if quote.IssuedAt != nil {
		q.IssuedAt = quote.IssuedAt
	}

	if quote.Plan.ID != 0 {
		q.Plan = new(Plan).Set(&quote.Plan)
	}

	if quote.RecipientAddress.ID != 0 {
		q.RecipientAddress = new(Address).Set(&quote.RecipientAddress)
	}

	if quote.User != nil {
		q.User = new(User).Set(quote.User)
	}

	return q
}
