package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type PlanFeature struct {
	ID        string    `json:"id"`
	Interval  string    `json:"interval"`
	Quota     int64     `json:"quota"`
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
