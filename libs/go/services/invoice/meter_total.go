package invoice

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func meterTotalForItem(db dbx.IORM, subscriptionID uint, si *models.SubscriptionItem) (int64, bool) {
	if si.PlanItemID == nil || !isMeteredPricing(si.PricingModel) {
		return 0, false
	}
	var meters []models.BillingMeter
	if err := db.List(&meters, "subscription_id = ? AND plan_item_id = ?", subscriptionID, *si.PlanItemID); err != nil {
		return 0, false
	}
	if len(meters) == 0 {
		return 0, false
	}
	total := int64(0)
	for i := range meters {
		billable := applyFreeQuota(meters[i].Usage, meters[i].FreeQuota)
		total += calculateItemTotal(
			meters[i].PricingModel,
			meters[i].UnitAmount,
			si.Quantity,
			billable,
			meters[i].Tiers,
		)
	}
	return total, true
}

func isMeteredPricing(model string) bool {
	switch model {
	case "unit", "tiered", "block", "volume":
		return true
	default:
		return false
	}
}

func applyFreeQuota(usage, quota int64) int64 {
	if usage <= 0 || quota <= 0 {
		if usage < 0 {
			return 0
		}
		return usage
	}
	if usage <= quota {
		return 0
	}
	return usage - quota
}
