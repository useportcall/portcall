package invoice

import "github.com/useportcall/portcall/libs/go/dbx/models"

// calculateItemTotal calculates the total for an invoice item based on pricing model.
func calculateItemTotal(
	pricingModel string,
	unitAmount int64,
	quantity int32,
	usage int64,
	tiers *[]models.Tier,
) int64 {
	switch pricingModel {
	case "fixed":
		return unitAmount * int64(quantity)
	case "unit":
		return unitAmount * usage * int64(quantity)
	case "tiered":
		return tieredTotal(tiers, usage) * int64(quantity)
	case "block":
		return blockTotal(tiers, usage) * int64(quantity)
	case "volume":
		return volumeTotal(tiers, usage) * int64(quantity)
	default:
		return 0
	}
}

func tieredTotal(tiers *[]models.Tier, usage int64) int64 {
	return graduatedTierTotal(tiers, usage)
}

func blockTotal(tiers *[]models.Tier, usage int64) int64 {
	return blockTierTotal(tiers, usage)
}

func volumeTotal(tiers *[]models.Tier, usage int64) int64 {
	return volumeTierTotal(tiers, usage)
}

// getItemUnitAmount returns the display unit amount for a subscription item.
func getItemUnitAmount(si models.SubscriptionItem) int64 {
	switch si.PricingModel {
	case "fixed", "unit":
		return si.UnitAmount
	default:
		if si.Tiers != nil && len(*si.Tiers) > 0 {
			return int64((*si.Tiers)[0].Amount)
		}
		return 0
	}
}

// calculateDiscount calculates the discount amount based on percentage.
func calculateDiscount(amount int64, discountPct int) int64 {
	if discountPct <= 0 {
		return 0
	}
	return amount * int64(discountPct) / 100
}
