package utils

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// calculateBlockTotal computes block pricing.
// Returns the flat fee for the tier that contains the usage.
func calculateBlockTotal(usage int64, tiers *[]models.Tier) int64 {
	if tiers == nil || len(*tiers) == 0 {
		return 0
	}

	for _, tier := range *tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		amount := int64(tier.Amount) // Flat amount for this tier

		if end == -1 {
			// Unlimited tier
			if usage > start {
				return amount
			}
		} else {
			// Bounded tier
			if usage > start && usage <= end {
				return amount
			}
		}
	}

	return 0
}

// calculateVolumeTotal computes volume pricing.
// The tier rate applies to ALL units, not just units within the tier.
func calculateVolumeTotal(usage int64, tiers *[]models.Tier) int64 {
	if tiers == nil || len(*tiers) == 0 {
		return 0
	}

	for _, tier := range *tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		amount := int64(tier.Amount) // Per-unit amount for this tier

		if end == -1 {
			// Unlimited tier
			if usage > start {
				return usage * amount
			}
		} else {
			// Bounded tier
			if usage > start && usage <= end {
				return usage * amount
			}
		}
	}

	return 0
}

// CalculateBillableUsage determines billable usage accounting for free quotas.
// Free quota is the amount that doesn't get charged (e.g., first 1000 free).
func CalculateBillableUsage(totalUsage int64, freeQuota int64) int64 {
	if freeQuota <= 0 {
		return totalUsage
	}
	if totalUsage <= freeQuota {
		return 0
	}
	return totalUsage - freeQuota
}
