package entitlement

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func syncMeteredUsage(db dbx.IORM, evt *models.MeterEvent, feature *models.Feature) error {
	if !feature.IsMetered {
		return nil
	}
	sub, err := findActiveSubscription(db, evt)
	if err != nil || sub == nil || sub.PlanID == nil {
		return err
	}

	var pf models.PlanFeature
	if err := db.FindFirst(&pf, "plan_id = ? AND feature_id = ?", *sub.PlanID, feature.ID); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			return nil
		}
		return err
	}

	if err := incrementSubItemUsage(db, sub.ID, pf.PlanItemID, evt.Usage); err != nil {
		return err
	}
	return upsertBillingMeter(db, sub, &pf, feature.ID, evt)
}

func findActiveSubscription(db dbx.IORM, evt *models.MeterEvent) (*models.Subscription, error) {
	var sub models.Subscription
	err := db.FindFirst(&sub,
		"app_id = ? AND user_id = ? AND status IN (?, ?)",
		evt.AppID, evt.UserID, "active", "resetting",
	)
	if dbx.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func incrementSubItemUsage(db dbx.IORM, subID, planItemID uint, delta int64) error {
	var item models.SubscriptionItem
	if err := db.FindFirst(&item, "subscription_id = ? AND plan_item_id = ?", subID, planItemID); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			return nil
		}
		return err
	}
	item.Usage = uint(clampUsage(int64(item.Usage), delta))
	return db.Save(&item)
}

func clampUsage(current, delta int64) int64 {
	next := current + delta
	if next < 0 {
		return 0
	}
	return next
}
