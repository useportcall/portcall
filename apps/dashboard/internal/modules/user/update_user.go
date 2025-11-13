package user

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func UpdateUser(c *routerx.Context) {
	id := c.Param("id")

	var body UpdateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), id, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	user.Name = body.Name
	if err := c.DB().Save(&user); err != nil {
		c.ServerError("Failed to save user", err)
		return
	}

	if user.BillingAddressID != nil {
		var billingAddress models.Address
		if err := c.DB().FindForID(*user.BillingAddressID, &billingAddress); err != nil {
			c.ServerError("Failed to fetch billing address", err)
			return
		}
		user.BillingAddress = &billingAddress
	}

	c.OK(new(User).Set(&user))
}
