package subscription

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListSubscriptions(c *routerx.Context) {
	subscriptions := []models.Subscription{}
	if err := c.DB().List(&subscriptions, "app_id = ?", c.AppID()); err != nil {
		c.ServerError("Failed to list subscriptions", err)
		return
	}

	result := make([]apix.Subscription, len(subscriptions))

	for i, sub := range subscriptions {
		s := new(apix.Subscription).Set(&sub)

		var u models.User
		if err := c.DB().FindForID(sub.UserID, &u); err == nil {
			s.User = new(apix.User).Set(&u)
		}

		if sub.PlanID != nil {
			var p models.Plan
			if err := c.DB().FindForID(*sub.PlanID, &p); err == nil {
				s.Plan = new(apix.Plan).Set(&p)
			}
		}

		var subscriptionItems []models.SubscriptionItem
		if err := c.DB().List(&subscriptionItems, "subscription_id = ?", sub.ID); err != nil {
			c.ServerError("Failed to list subscription items", err)
			return
		}

		s.Items = make([]apix.SubscriptionItem, len(subscriptionItems))
		for i, item := range subscriptionItems {
			s.Items[i].Set(&item)
		}

		result[i] = *s
	}

	c.OK(result)
}
