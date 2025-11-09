package company

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/address"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetCompany(c *routerx.Context) {
	company := models.Company{}
	if err := c.DB().FindFirstForAppID(c.AppID(), &company); err != nil {
		c.NotFound("Company not found")
		return
	}

	response := new(Company).Set(&company)

	var billingAddress models.Address
	if err := c.DB().FindForID(company.BillingAddressID, &billingAddress); err == nil {
		response.BillingAddress = new(address.Address)
		response.BillingAddress.Set(&billingAddress)
	}

	c.OK(response)
}
