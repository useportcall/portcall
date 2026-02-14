package utils

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// CalculateBillingTotal calculates the total billing amount for metered usage.
// This properly handles tiered, block, and unit pricing models.
//
// For tiered pricing: Each tier charges for usage within that tier's range.
//
//	Example: 0-1000 free, 1001-5000 at $0.10, 5001+ at $0.05
//	If usage is 6000: (0 * 1000) + (0.10 * 4000) + (0.05 * 1000) = $450.00
//
// For block pricing: A flat fee based on which tier the total usage falls into.
//
//	Example: 0-1000 = $10, 1001-5000 = $40, 5001+ = $100
//	If usage is 6000: $100 flat
//
// For unit pricing: Simple multiplication of usage * unit amount.
func CalculateBillingTotal(
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
		return calculateTieredTotal(usage, tiers)

	case "block":
		return calculateBlockTotal(usage, tiers)

	case "volume":
		// Volume pricing applies the tier's rate to ALL units
		return calculateVolumeTotal(usage, tiers)

	default:
		return 0
	}
}

// calculateTieredTotal computes graduated tiered pricing.
// Each tier is charged only for usage within its range.
func calculateTieredTotal(usage int64, tiers *[]models.Tier) int64 {
	if tiers == nil || len(*tiers) == 0 {
		return 0
	}

	var total int64 = 0

	for _, tier := range *tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		amount := int64(tier.Amount) // Amount in cents per unit

		// Calculate units in this tier
		var tierUsage int64
		if end == -1 {
			// Unlimited tier - all remaining usage
			if usage > start {
				tierUsage = usage - start
			}
		} else {
			// Bounded tier
			if usage <= start {
				continue
			}
			tierMax := end - start
			if usage >= end {
				tierUsage = tierMax
			} else {
				tierUsage = usage - start
			}
		}

		total += tierUsage * amount
	}

	return total
}
