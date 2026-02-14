package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// BillingMeter represents billing usage tracking for metered features.
// This tracks usage over the subscription's billing period for invoicing.
type BillingMeter struct {
	ID             string     `json:"id"`
	SubscriptionID string     `json:"subscription_id"`
	FeatureID      string     `json:"feature_id"`
	PlanItemID     string     `json:"plan_item_id"`
	UserID         string     `json:"user_id"`
	Usage          int64      `json:"usage"`
	PricingModel   string     `json:"pricing_model"`
	UnitAmount     int64      `json:"unit_amount"`
	FreeQuota      int64      `json:"free_quota"`
	Tiers          []Tier     `json:"tiers,omitempty"`
	Interval       string     `json:"interval"`       // billing interval: inherit, week, month, year
	IntervalCount  int        `json:"interval_count"` // number of intervals for the billing cycle
	LastResetAt    *time.Time `json:"last_reset_at"`
	NextResetAt    *time.Time `json:"next_reset_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Computed fields for display
	ProjectedCost int64  `json:"projected_cost"` // Estimated cost at current usage
	FeatureName   string `json:"feature_name,omitempty"`
	PlanItemTitle string `json:"plan_item_title,omitempty"`
}

// Tier represents a pricing tier for tiered/block/volume pricing.
type Tier struct {
	Start  int `json:"start"`
	End    int `json:"end"`
	Amount int `json:"unit_amount"`
}

func (b *BillingMeter) Set(meter *models.BillingMeter) *BillingMeter {
	b.Usage = meter.Usage
	b.PricingModel = meter.PricingModel
	b.UnitAmount = meter.UnitAmount
	b.FreeQuota = meter.FreeQuota
	b.Interval = meter.Interval
	b.IntervalCount = meter.IntervalCount
	b.LastResetAt = meter.LastResetAt
	b.NextResetAt = meter.NextResetAt
	b.CreatedAt = meter.CreatedAt
	b.UpdatedAt = meter.UpdatedAt

	// Convert tiers
	if meter.Tiers != nil && len(*meter.Tiers) > 0 {
		b.Tiers = make([]Tier, len(*meter.Tiers))
		for i, t := range *meter.Tiers {
			b.Tiers[i] = Tier{
				Start:  t.Start,
				End:    t.End,
				Amount: t.Amount,
			}
		}
	}

	return b
}

// CalculateProjectedCost calculates the projected cost at current usage.
func (b *BillingMeter) CalculateProjectedCost() int64 {
	switch b.PricingModel {
	case "unit":
		return b.UnitAmount * b.Usage
	case "tiered":
		return b.calculateTieredCost()
	case "block":
		return b.calculateBlockCost()
	case "volume":
		return b.calculateVolumeCost()
	default:
		return 0
	}
}

func (b *BillingMeter) calculateTieredCost() int64 {
	if len(b.Tiers) == 0 {
		return 0
	}

	var total int64 = 0
	usage := b.Usage

	for _, tier := range b.Tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		amount := int64(tier.Amount)

		var tierUsage int64
		if end == -1 {
			if usage > start {
				tierUsage = usage - start
			}
		} else {
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

func (b *BillingMeter) calculateBlockCost() int64 {
	if len(b.Tiers) == 0 {
		return 0
	}

	for _, tier := range b.Tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		amount := int64(tier.Amount)

		if end == -1 {
			if b.Usage > start {
				return amount
			}
		} else {
			if b.Usage > start && b.Usage <= end {
				return amount
			}
		}
	}

	return 0
}

func (b *BillingMeter) calculateVolumeCost() int64 {
	if len(b.Tiers) == 0 {
		return 0
	}

	for _, tier := range b.Tiers {
		start := int64(tier.Start)
		end := int64(tier.End)
		amount := int64(tier.Amount)

		if end == -1 {
			if b.Usage > start {
				return b.Usage * amount
			}
		} else {
			if b.Usage > start && b.Usage <= end {
				return b.Usage * amount
			}
		}
	}

	return 0
}
