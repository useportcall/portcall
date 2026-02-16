package users

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UserListItem struct {
	ID               uint      `json:"id"`
	PublicID         string    `json:"public_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	HasSubscription  bool      `json:"has_subscription"`
	HasPaymentMethod bool      `json:"has_payment_method"`
	EntitlementCount int64     `json:"entitlement_count"`
	CreatedAt        time.Time `json:"created_at"`
}

func ListUsers(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid app ID")
		return
	}

	var users []models.User
	if err := c.DB().List(&users, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list users", err)
		return
	}

	result := make([]UserListItem, len(users))
	for i, user := range users {
		var subCount, entCount, pmCount int64
		c.DB().Count(&subCount, models.Subscription{}, "user_id = ?", user.ID)
		c.DB().Count(&entCount, models.Entitlement{}, "user_id = ?", user.ID)
		c.DB().Count(&pmCount, models.PaymentMethod{}, "user_id = ?", user.ID)

		result[i] = UserListItem{
			ID:               user.ID,
			PublicID:         user.PublicID,
			Name:             user.Name,
			Email:            user.Email,
			HasSubscription:  subCount > 0,
			HasPaymentMethod: pmCount > 0,
			EntitlementCount: entCount,
			CreatedAt:        user.CreatedAt,
		}
	}

	c.OK(result)
}
