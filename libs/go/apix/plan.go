package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Plan struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Currency        string     `json:"currency"`
	Status          string     `json:"status"`
	TrialPeriodDays int        `json:"trial_period_days"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Items           []PlanItem `json:"items"`
	Interval        string     `json:"interval"`
	IntervalCount   int        `json:"interval_count"`
	PlanGroup       any        `json:"plan_group"`
	Features        []any      `json:"features"`
	MeteredFeatures []any      `json:"metered_features"`
}

func (p *Plan) Set(plan *models.Plan) *Plan {
	p.ID = plan.PublicID
	p.Name = plan.Name
	p.Currency = plan.Currency
	p.Status = plan.Status
	p.TrialPeriodDays = plan.TrialPeriodDays
	p.CreatedAt = plan.CreatedAt
	p.UpdatedAt = plan.UpdatedAt
	p.Interval = plan.Interval
	p.IntervalCount = plan.IntervalCount
	return p
}
