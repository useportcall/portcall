package subscription

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateSubscriptionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	PlanID string `json:"plan_id" binding:"required"`
}

func CreateSubscription(c *routerx.Context) {
	if err := checkMaxSubscriptions(c); err != nil {
		return
	}
	var body CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("invalid request body: user_id and plan_id are required")
		return
	}
	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), body.UserID, &user); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("User not found")
			return
		}
		c.ServerError("error retrieving user", err)
		return
	}
	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), body.PlanID, &plan); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("Plan not found")
			return
		}
		c.ServerError("error retrieving plan", err)
		return
	}
	if !plan.IsFree {
		hasPM, err := hasPaymentMethod(c, user.ID)
		if err != nil {
			c.ServerError("error checking payment method", err)
			return
		}
		if !hasPM {
			c.BadRequest("paid plan requires a saved payment method; use checkout or payment link")
			return
		}
	}

	createdNew, err := runBillingFlow(c, user.ID, plan.ID)
	if err != nil {
		c.ServerError("error applying subscription flow", err)
		return
	}
	if createdNew {
		recordDogfoodUsage(c, &user)
	}

	var sub models.Subscription
	if err := c.DB().FindFirst(&sub, "app_id = ? AND user_id = ? AND status = ?", c.AppID(), user.ID, "active"); err != nil {
		c.ServerError("error loading subscription", err)
		return
	}
	sub.User = user
	sub.Plan = &plan
	c.OK(new(apix.Subscription).Set(&sub))
}
