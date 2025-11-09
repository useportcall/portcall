package address

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func UpdateAddress(c *routerx.Context) {
	id := c.Param("id")

	var body UpdateAddressRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var address models.Address
	if err := c.DB().GetForPublicID(c.AppID(), id, &address); err != nil {
		c.NotFound("Address not found")
		return
	}

	if body.Line1 != "" {
		address.Line1 = body.Line1
	}

	if body.Line2 != "" {
		address.Line2 = body.Line2
	}

	if body.City != "" {
		address.City = body.City
	}

	if body.State != "" {
		address.State = body.State
	}

	if body.PostalCode != "" {
		address.PostalCode = body.PostalCode
	}

	if body.Country != "" {
		address.Country = body.Country
	}

	if err := c.DB().Save(&address); err != nil {
		c.ServerError("Failed to update address")
		return
	}

	c.OK(new(Address).Set(&address))
}
