package entitlement

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func upsertBillingMeter(
	db dbx.IORM, sub *models.Subscription, pf *models.PlanFeature,
	featureID uint, evt *models.MeterEvent,
) error {
	var meter models.BillingMeter
	err := db.FindFirst(&meter,
		"subscription_id = ? AND plan_item_id = ? AND feature_id = ?",
		sub.ID, pf.PlanItemID, featureID,
	)
	if dbx.IsRecordNotFoundError(err) {
		return createBillingMeter(db, sub, pf, featureID, evt)
	}
	if err != nil {
		return err
	}

	meter.Usage = clampUsage(meter.Usage, evt.Usage)
	return db.Save(&meter)
}

func createBillingMeter(
	db dbx.IORM, sub *models.Subscription, pf *models.PlanFeature,
	featureID uint, evt *models.MeterEvent,
) error {
	var planItem models.PlanItem
	if err := db.FindForID(pf.PlanItemID, &planItem); err != nil {
		return err
	}

	interval, count := planItem.Interval, planItem.IntervalCount
	if interval == "" || interval == "inherit" {
		interval = sub.BillingInterval
		count = sub.BillingIntervalCount
	}

	freeQuota := int64(0)
	if pf.Quota > 0 {
		freeQuota = pf.Quota
	}

	meter := models.BillingMeter{
		AppID:          sub.AppID,
		SubscriptionID: sub.ID,
		PlanItemID:     planItem.ID,
		FeatureID:      featureID,
		UserID:         evt.UserID,
		Usage:          clampUsage(0, evt.Usage),
		PricingModel:   planItem.PricingModel,
		UnitAmount:     planItem.UnitAmount,
		Tiers:          planItem.Tiers,
		FreeQuota:      freeQuota,
		Interval:       interval,
		IntervalCount:  count,
		LastResetAt:    &sub.LastResetAt,
		NextResetAt:    &sub.NextResetAt,
	}
	return db.Create(&meter)
}
