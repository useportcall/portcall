package address

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateAddress(c *routerx.Context) {
	var body CreateAddressRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	address := new(models.Address)
	address.PublicID = utils.GenPublicID("addr")
	address.AppID = c.AppID()
	address.Line1 = body.Line1
	address.Line2 = body.Line2
	address.City = body.City
	address.State = body.State
	address.PostalCode = body.PostalCode
	address.Country = body.Country

	if err := c.DB().Create(address); err != nil {
		c.ServerError("Failed to create address")
		return
	}

	var company models.Company
	if err := c.DB().FindFirstForAppID(c.AppID(), &company); err != nil {
		c.ServerError("Failed to find company")
		return
	}

	company.BillingAddressID = address.ID
	if err := c.DB().Save(&company); err != nil {
		c.ServerError("Failed to update company with billing address")
		return
	}

	c.OK(new(Address).Set(address))
}
