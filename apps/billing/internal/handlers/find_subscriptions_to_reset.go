package handlers

import (
	"log"

	"github.com/useportcall/portcall/libs/go/qx/server"
)

func FindSubscriptionsToReset(c server.IContext) error {
	subscriptionIDs := []uint{}
	if err := c.DB().ListIDs(
		"subscriptions",
		&subscriptionIDs,
		"(status = ? OR status = ?) AND next_reset_at < CURRENT_TIMESTAMP",
		"active",
		"canceled",
	); err != nil {
		return err
	}

	if len(subscriptionIDs) == 0 {
		log.Println("No subscriptions to reset found")
		return nil
	}

	for _, id := range subscriptionIDs {
		if err := c.Queue().Enqueue("start_subscription_reset", map[string]any{"subscription_id": id}, "billing_queue"); err != nil {
			return err
		}
	}

	return nil
}
