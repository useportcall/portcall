package user

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetUser(c *routerx.Context) {
	id := c.Param("id")

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), id, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	if user.BillingAddressID != nil {
		var billingAddress models.Address
		if err := c.DB().FindForID(*user.BillingAddressID, &billingAddress); err != nil {
			c.ServerError("Failed to get billing address")
			return
		}
		user.BillingAddress = &billingAddress
	}

	c.OK(new(User).Set(&user))
}
