package user

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListUsers(c *routerx.Context) {
	email := c.Query("email")
	var conds []any
	if email != "" {
		conds = append(conds, "app_id = ? AND email LIKE ?", c.AppID(), "%"+email+"%")
	} else {
		conds = append(conds, "app_id = ?", c.AppID())
	}

	var users []models.User
	if err := c.DB().List(&users, conds...); err != nil {
		c.ServerError("Failed to list users", err)
		return
	}

	response := make([]apix.User, len(users))
	for i := range users {
		if users[i].BillingAddressID != nil {
			var billingAddress models.Address
			if err := c.DB().FindForID(*users[i].BillingAddressID, &billingAddress); err == nil {
				users[i].BillingAddress = &billingAddress
			}
		}

		response[i] = *new(apix.User).Set(&users[i])

		var subscriptionCount int64
		if err := c.DB().Count(&subscriptionCount, models.Subscription{}, "user_id = ?", users[i].ID); err == nil {
			response[i].Subscribed = subscriptionCount > 0
		}

		var paymentMethodCount int64
		if err := c.DB().Count(&paymentMethodCount, models.PaymentMethod{}, "user_id = ?", users[i].ID); err == nil {
			response[i].PaymentMethodAdded = paymentMethodCount > 0
		}
	}

	c.OK(response)
}
