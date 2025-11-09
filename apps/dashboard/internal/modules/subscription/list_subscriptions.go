package subscription

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/subscription_item"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/user"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListSubscriptions(c *routerx.Context) {
	subscriptions := []models.Subscription{}
	if err := c.DB().List(&subscriptions, "app_id = ?", c.AppID()); err != nil {
		c.ServerError("Failed to list subscriptions")
		return
	}

	result := make([]Subscription, len(subscriptions))

	for i, sub := range subscriptions {
		s := new(Subscription).Set(&sub)

		var u models.User
		if err := c.DB().FindForID(sub.UserID, &u); err == nil {
			s.User = new(user.User).Set(&u)
		}

		if sub.PlanID != nil {
			var p models.Plan
			if err := c.DB().FindForID(*sub.PlanID, &p); err == nil {
				s.Plan = new(plan.Plan).Set(&p)
			}
		}

		var subscriptionItems []models.SubscriptionItem
		if err := c.DB().List(&subscriptionItems, "subscription_id = ?", sub.ID); err == nil {
			for _, item := range subscriptionItems {
				s.Items = append(s.Items, *new(subscription_item.SubscriptionItem).Set(&item))
			}
		}

		result[i] = *s
	}

	c.OK(result)
}
