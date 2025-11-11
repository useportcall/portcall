package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Address struct {
	ID         string    `json:"id"`
	Line1      string    `json:"line1"`
	Line2      string    `json:"line2,omitempty"`
	City       string    `json:"city"`
	State      string    `json:"state,omitempty"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (a *Address) Set(address *models.Address) *Address {
	a.ID = address.PublicID
	a.Line1 = address.Line1
	a.Line2 = address.Line2
	a.City = address.City
	a.State = address.State
	a.PostalCode = address.PostalCode
	a.Country = address.Country
	a.CreatedAt = address.CreatedAt
	a.UpdatedAt = address.UpdatedAt
	return a
}
