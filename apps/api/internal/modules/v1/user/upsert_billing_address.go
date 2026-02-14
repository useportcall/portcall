package user

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpsertBillingAddressRequest struct {
	Line1      string `json:"line1" binding:"required"`
	Line2      string `json:"line2"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code" binding:"required"`
	Country    string `json:"country" binding:"required"`
}

func UpsertBillingAddress(c *routerx.Context) {
	id := c.Param("id")

	var body UpsertBillingAddressRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), id, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	var address models.Address

	// If user already has a billing address, update it
	if user.BillingAddressID != nil {
		if err := c.DB().FindForID(*user.BillingAddressID, &address); err != nil {
			c.ServerError("Failed to fetch billing address", err)
			return
		}

		address.Line1 = body.Line1
		address.Line2 = body.Line2
		address.City = body.City
		address.State = body.State
		address.PostalCode = body.PostalCode
		address.Country = body.Country

		if err := c.DB().Save(&address); err != nil {
			c.ServerError("Failed to update billing address", err)
			return
		}
	} else {
		// Create a new billing address
		address = models.Address{
			PublicID:   dbx.GenPublicID("addr"),
			AppID:      c.AppID(),
			Line1:      body.Line1,
			Line2:      body.Line2,
			City:       body.City,
			State:      body.State,
			PostalCode: body.PostalCode,
			Country:    body.Country,
		}

		if err := c.DB().Create(&address); err != nil {
			c.ServerError("Failed to create billing address", err)
			return
		}

		// Update user with billing address ID
		user.BillingAddressID = &address.ID
		if err := c.DB().Save(&user); err != nil {
			c.ServerError("Failed to update user with billing address", err)
			return
		}
	}

	c.OK(new(apix.Address).Set(&address))
}
