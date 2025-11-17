package company

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetCompany(c *routerx.Context) {
	company := models.Company{}
	if err := c.DB().FindFirstForAppID(c.AppID(), &company); err != nil {
		c.NotFound("Company not found")
		return
	}

	response := new(apix.Company).Set(&company)

	var billingAddress models.Address
	if err := c.DB().FindForID(company.BillingAddressID, &billingAddress); err == nil {
		response.BillingAddress = new(apix.Address)
		response.BillingAddress.Set(&billingAddress)
	}

	c.OK(response)
}
