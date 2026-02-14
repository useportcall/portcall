package invoice

import "github.com/useportcall/portcall/libs/go/dbx/models"

func graduatedTierTotal(tiers *[]models.Tier, usage int64) int64 {
	if tiers == nil || len(*tiers) == 0 || usage <= 0 {
		return 0
	}
	total := int64(0)
	for _, tier := range *tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		if usage <= start {
			continue
		}
		upper := usage
		if end != -1 && usage > end {
			upper = end
		}
		units := upper - start
		if units > 0 {
			total += units * int64(tier.Amount)
		}
	}
	return total
}

func blockTierTotal(tiers *[]models.Tier, usage int64) int64 {
	if tiers == nil || len(*tiers) == 0 || usage <= 0 {
		return 0
	}
	for _, tier := range *tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		if usage > start && (end == -1 || usage <= end) {
			return int64(tier.Amount)
		}
	}
	return 0
}

func volumeTierTotal(tiers *[]models.Tier, usage int64) int64 {
	if tiers == nil || len(*tiers) == 0 || usage <= 0 {
		return 0
	}
	for _, tier := range *tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		if usage > start && (end == -1 || usage <= end) {
			return usage * int64(tier.Amount)
		}
	}
	return 0
}
