package company

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func UpsertCompany(c *routerx.Context) {
	var body UpdateCompanyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	company := models.Company{}
	if err := c.DB().FindFirstForAppID(c.AppID(), &company); err != nil {
		c.NotFound("Company not found")
		return
	}

	if body.Name != nil {
		company.Name = *body.Name
	}

	if body.FirstName != nil {
		company.FirstName = *body.FirstName
	}

	if body.LastName != nil {
		company.LastName = *body.LastName
	}

	if body.Email != nil {
		company.Email = *body.Email
	}

	if body.Phone != nil {
		company.Phone = *body.Phone
	}

	if body.VATNumber != nil {
		company.VATNumber = *body.VATNumber
	}

	if body.Alias != nil {
		company.Alias = *body.Alias
	}

	if body.BusinessCategory != nil {
		company.BusinessCategory = *body.BusinessCategory
	}

	if err := c.DB().Save(&company); err != nil {
		c.ServerError("Failed to update company", err)
		return
	}

	c.OK(new(apix.Company).Set(&company))
}

type UpdateCompanyRequest struct {
	Name             *string `json:"name"`
	Alias            *string `json:"alias"`
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	Phone            *string `json:"phone"`
	VATNumber        *string `json:"vat_number"`
	BusinessCategory *string `json:"business_category"`
}
