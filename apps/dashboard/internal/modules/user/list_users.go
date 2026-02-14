package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListUsers(c *routerx.Context) {
	search := c.Query("q")
	if search == "" {
		search = c.Query("email")
	}
	query := "app_id = ?"
	args := []any{c.AppID()}
	query, args = appendUserSearchFilter(query, args, search)
	query = appendUserExistsFilter(query, c.Query("subscribed"), "subscriptions")
	query = appendUserExistsFilter(query, c.Query("payment_method_added"), "payment_methods")

	rawLimit := c.Query("limit")
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(rawLimit)
	isPaginated := rawLimit != ""
	if limit < 1 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}

	var users []models.User
	if isPaginated {
		offset := (page - 1) * limit
		if err := c.DB().ListWithOrderLimitOffset(&users, "created_at DESC", limit, offset, append([]any{query}, args...)...); err != nil {
			c.ServerError("Failed to list users", err)
			return
		}
	} else if err := c.DB().List(&users, append([]any{query}, args...)...); err != nil {
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

	if !isPaginated {
		c.OK(response)
		return
	}

	var total int64
	if err := c.DB().Count(&total, models.User{}, query, args...); err != nil {
		c.ServerError("Failed to count users", err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	c.OK(gin.H{
		"users":       response,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	})
}
