package plan_feature

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type PlanFeature struct {
	ID        string    `json:"id"`
	Interval  string    `json:"interval"`
	Quota     int       `json:"quota"`
	Rollover  int       `json:"rollover"`
	Feature   any       `json:"feature"`
	PlanItem  any       `json:"plan_item"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (pf *PlanFeature) Set(planFeature *models.PlanFeature) *PlanFeature {
	pf.ID = planFeature.PublicID
	pf.Interval = planFeature.Interval
	pf.Quota = planFeature.Quota
	pf.Rollover = planFeature.Rollover
	pf.CreatedAt = planFeature.CreatedAt
	pf.UpdatedAt = planFeature.UpdatedAt
	return pf
}

type UpdatePlanFeatureRequest struct {
	PlanItemID string  `json:"plan_item_id"`
	FeatureID  *string `json:"feature_id"`
	Interval   string  `json:"interval"`
	Quota      int     `json:"quota"`
	Rollover   *int    `json:"rollover"`
}

type CreatePlanFeatureRequest struct {
	FeatureID string `json:"feature_id"`
	PlanID    string `json:"plan_id"`
}
