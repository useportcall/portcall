package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type PlanItem struct {
	ID                string         `json:"id"`
	Quantity          int32          `json:"quantity"`
	PricingModel      string         `json:"pricing_model"`
	UnitAmount        int64          `json:"unit_amount"`
	Tiers             *[]models.Tier `json:"tiers"`
	Minimum           *int           `json:"minimum"` //TODO: implement
	Maximum           *int           `json:"maximum"` //TODO: implement
	PublicTitle       string         `json:"public_title"`
	PublicDescription string         `json:"public_description"`
	PublicUnitLabel   string         `json:"public_unit_label"`
	Interval          string         `json:"interval"`       // billing interval: inherit, week, month, year
	IntervalCount     int            `json:"interval_count"` // number of intervals for the billing cycle
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`

	Features []any `json:"features"`
}

func (pi *PlanItem) Set(planItem *models.PlanItem) *PlanItem {
	pi.ID = planItem.PublicID
	pi.Quantity = planItem.Quantity
	pi.PricingModel = planItem.PricingModel
	pi.UnitAmount = planItem.UnitAmount
	pi.Tiers = planItem.Tiers
	pi.Minimum = planItem.Minimum
	pi.Maximum = planItem.Maximum
	pi.PublicTitle = planItem.PublicTitle
	pi.PublicDescription = planItem.PublicDescription
	pi.PublicUnitLabel = planItem.PublicUnitLabel
	pi.Interval = planItem.Interval
	pi.IntervalCount = planItem.IntervalCount
	pi.CreatedAt = planItem.CreatedAt
	pi.UpdatedAt = planItem.UpdatedAt
	return pi
}
