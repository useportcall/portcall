package plan_item

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
	pi.CreatedAt = planItem.CreatedAt
	pi.UpdatedAt = planItem.UpdatedAt
	return pi
}

type CreatePlanItemRequest struct {
	PlanID            string `json:"plan_id"`
	PricingModel      string `json:"pricing_model"`
	UnitAmount        int64  `json:"unit_amount"`
	PublicTitle       string `json:"public_title"`
	PublicDescription string `json:"public_description"`
	Interval          string `json:"interval"`
	Quota             int    `json:"quota"`
	Rollover          int    `json:"rollover"`
}

type UpdatePlanItemRequest struct {
	PricingModel      string        `json:"pricing_model"`
	Quantity          int32         `json:"quantity"`
	UnitAmount        int64         `json:"unit_amount"`
	Tiers             []models.Tier `json:"tiers"`
	Minimum           *int          `json:"minimum"`
	Maximum           *int          `json:"maximum"`
	PublicTitle       string        `json:"public_title"`
	PublicDescription string        `json:"public_description"`
}
