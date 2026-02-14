package user

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateUserRequest struct {
	Name         *string `json:"name"`
	CompanyTitle *string `json:"company_title"`
}

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

	if body.Name != nil {
		user.Name = *body.Name
	}

	if body.CompanyTitle != nil {
		user.CompanyTitle = *body.CompanyTitle
	}

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

	c.OK(new(apix.User).Set(&user))
}
