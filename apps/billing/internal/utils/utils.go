package utils

import "github.com/useportcall/portcall/libs/go/dbx/models"

func IsPricingMetered(pricingModel string) bool {
	switch pricingModel {
	case "fixed":
		return false
	default:
		return true
	}
}

func CalculateTotal(
	pricingModel string,
	unitAmount int64,
	quantity int32,
	usage uint,
	tiers *[]models.Tier,
) int64 {
	switch pricingModel {
	case "fixed":
		return int64(unitAmount) * int64(quantity)
	case "unit":
		return int64(unitAmount) * int64(usage) * int64(quantity)
	case "tiered":
		if tiers == nil {
			return 0
		}

		for _, tier := range *tiers {
			usage := int64(usage)
			start := int64(tier.Start)
			end := int64(tier.End)

			if end == -1 || usage > start && int64(usage) <= end {
				return int64(tier.Amount) * int64(quantity) * int64(usage)
			}
		}

		return 0
	case "block":
		if tiers == nil {
			return 0
		}

		for _, tier := range *tiers {
			usage := int64(usage)
			start := int64(tier.Start)
			end := int64(tier.End)

			if end == int64(-1) || usage > start && usage <= end {
				return int64(tier.Amount) * int64(quantity)
			}
		}

		return 0
	default:
		return 0
	}
}

func GetItemUnitAmount(si models.SubscriptionItem) int64 {
	switch si.PricingModel {
	case "fixed", "unit":
		return si.UnitAmount
	case "tiered":
		if si.Tiers != nil && len(*si.Tiers) > 0 {
			return int64((*si.Tiers)[0].Amount)
		}
		return 0
	case "block":
		if si.Tiers != nil && len(*si.Tiers) > 0 {
			return int64((*si.Tiers)[0].Amount)
		}
		return 0
	case "volume":
		if si.Tiers != nil && len(*si.Tiers) > 0 {
			return int64((*si.Tiers)[0].Amount)
		}
		return 0
	default:
		return 0
	}
}
