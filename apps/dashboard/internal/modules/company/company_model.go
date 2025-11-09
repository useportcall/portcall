package company

import (
	"time"

	"github.com/useportcall/portcall/apps/dashboard/internal/modules/address"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Company struct {
	Name             string           `json:"name"`
	Alias            string           `json:"alias"`
	FirstName        string           `json:"first_name"`
	LastName         string           `json:"last_name"`
	Email            string           `json:"email"`
	Phone            string           `json:"phone"`
	VATNumber        string           `json:"vat_number"`
	BusinessCategory string           `json:"business_category"`
	BillingAddress   *address.Address `json:"billing_address"`
	ShippingAddress  *address.Address `json:"shipping_address"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

func (c *Company) Set(company *models.Company) *Company {
	c.Name = company.Name
	c.Alias = company.Alias
	c.FirstName = company.FirstName
	c.LastName = company.LastName
	c.Email = company.Email
	c.Phone = company.Phone
	c.VATNumber = company.VATNumber
	c.BusinessCategory = company.BusinessCategory
	c.CreatedAt = company.CreatedAt
	c.UpdatedAt = company.UpdatedAt
	return c
}
