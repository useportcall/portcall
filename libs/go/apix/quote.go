package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Quote struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	Email     string     `json:"email"` // email of the user to whom the quote is sent
	DaysValid int        `json:"days_valid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	IssuedAt  *time.Time `json:"issued_at"`
	Plan      any        `json:"plan"`
	User      any        `json:"user"`
}

func (q *Quote) Set(quote *models.Quote) *Quote {
	q.ID = quote.PublicID
	q.Status = quote.Status
	q.Email = quote.User.Email
	q.DaysValid = quote.DaysValid
	q.CreatedAt = quote.CreatedAt
	q.UpdatedAt = quote.UpdatedAt

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
