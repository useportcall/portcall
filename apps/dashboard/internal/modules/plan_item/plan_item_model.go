package plan_item

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

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
	Quantity          *int32        `json:"quantity"`
	UnitAmount        *int64        `json:"unit_amount"`
	Tiers             []models.Tier `json:"tiers"`
	Minimum           *int          `json:"minimum"`
	Maximum           *int          `json:"maximum"`
	PublicTitle       string        `json:"public_title"`
	PublicDescription string        `json:"public_description"`
}
