package subscription

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func resetUsageRecords(db dbx.IORM, subscriptionID uint, now, next time.Time) error {
	if err := resetSubscriptionItems(db, subscriptionID, now, next); err != nil {
		return err
	}
	return resetBillingMeters(db, subscriptionID, now, next)
}

func resetSubscriptionItems(db dbx.IORM, subscriptionID uint, now, next time.Time) error {
	var items []models.SubscriptionItem
	if err := db.List(&items, "subscription_id = ?", subscriptionID); err != nil {
		return err
	}
	for i := range items {
		items[i].Usage = 0
		items[i].LastResetAt = &now
		items[i].NextResetAt = &next
		if err := db.Save(&items[i]); err != nil {
			return err
		}
	}
	return nil
}

func resetBillingMeters(db dbx.IORM, subscriptionID uint, now, next time.Time) error {
	var meters []models.BillingMeter
	if err := db.List(&meters, "subscription_id = ?", subscriptionID); err != nil {
		return err
	}
	for i := range meters {
		meters[i].Usage = 0
		meters[i].LastResetAt = &now
		meters[i].NextResetAt = &next
		if err := db.Save(&meters[i]); err != nil {
			return err
		}
	}
	return nil
}
