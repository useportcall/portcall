package address

import (
	quotemodule "github.com/useportcall/portcall/apps/dashboard/internal/modules/quote"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateAddressRequest struct {
	Line1      *string `json:"line1"`
	Line2      *string `json:"line2"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
}

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
	locked, err := quotemodule.HasLockedQuoteForAddress(c, address.ID)
	if err != nil {
		c.ServerError("Failed to validate quote state", err)
		return
	}
	if locked {
		c.BadRequest("Address cannot be edited after quote is issued")
		return
	}

	if body.Line1 != nil {
		address.Line1 = *body.Line1
	}

	if body.Line2 != nil {
		address.Line2 = *body.Line2
	}

	if body.City != nil {
		address.City = *body.City
	}

	if body.State != nil {
		address.State = *body.State
	}

	if body.PostalCode != nil {
		address.PostalCode = *body.PostalCode
	}

	if body.Country != nil {
		address.Country = *body.Country
	}

	if err := c.DB().Save(&address); err != nil {
		c.ServerError("Failed to update address", err)
		return
	}

	c.OK(new(apix.Address).Set(&address))
}
