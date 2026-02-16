package resources

import (
	"context"
	"fmt"
	"time"
)

// Plan represents a billing plan
type Plan struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Description     *string       `json:"description,omitempty"`
	Type            string        `json:"type,omitempty"` // free, pro, enterprise
	Status          string        `json:"status"`
	Currency        string        `json:"currency"`
	Interval        string        `json:"interval"`
	IntervalCount   int           `json:"interval_count"`
	TrialPeriodDays int           `json:"trial_period_days"`
	UnitAmount      int64         `json:"unit_amount"`
	IsFree          bool          `json:"is_free"`
	Items           []PlanItem    `json:"items,omitempty"`
	Features        []PlanFeature `json:"features,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// PlanItem represents an item within a plan
type PlanItem struct {
	ID                string    `json:"id"`
	Quantity          *int      `json:"quantity,omitempty"`
	PricingModel      *string   `json:"pricing_model,omitempty"`
	UnitAmount        *int64    `json:"unit_amount,omitempty"`
	PublicTitle       *string   `json:"public_title,omitempty"`
	PublicDescription *string   `json:"public_description,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// PlanFeature represents a feature attached to a plan
type PlanFeature struct {
	ID        string    `json:"id"`
	FeatureID string    `json:"feature_id"`
	Feature   *string   `json:"feature,omitempty"`
	Interval  *string   `json:"interval,omitempty"`
	Quota     *int64    `json:"quota,omitempty"`
	Rollover  *int      `json:"rollover,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePlanRequest is the request to create a plan
type CreatePlanRequest struct {
	Name            string `json:"name"`
	Currency        string `json:"currency,omitempty"`
	Interval        string `json:"interval,omitempty"`
	IntervalCount   int    `json:"interval_count,omitempty"`
	TrialPeriodDays int    `json:"trial_period_days,omitempty"`
	UnitAmount      int64  `json:"unit_amount,omitempty"`
}

// UpdatePlanRequest is the request to update a plan
type UpdatePlanRequest struct {
	Name            *string `json:"name,omitempty"`
	Currency        *string `json:"currency,omitempty"`
	Interval        *string `json:"interval,omitempty"`
	IntervalCount   *int    `json:"interval_count,omitempty"`
	TrialPeriodDays *int    `json:"trial_period_days,omitempty"`
	UnitAmount      *int64  `json:"unit_amount,omitempty"`
}

// CreatePlanFeatureRequest is the request to add a feature to a plan
type CreatePlanFeatureRequest struct {
	PlanID    string `json:"plan_id"`
	FeatureID string `json:"feature_id"`
	Interval  string `json:"interval,omitempty"`
	Quota     *int64 `json:"quota,omitempty"`
}

// Plans provides access to plan-related API operations
type Plans struct {
	http *HTTPClient
}

// NewPlans creates a new Plans resource
func NewPlans(http *HTTPClient) *Plans {
	return &Plans{http: http}
}

// List returns all plans for the app
func (p *Plans) List(ctx context.Context) ([]Plan, error) {
	var resp DataWrapper[[]Plan]
	if err := p.http.Get(ctx, "/v1/plans", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a plan by ID
func (p *Plans) Get(ctx context.Context, planID string) (*Plan, error) {
	var resp DataWrapper[Plan]
	if err := p.http.Get(ctx, fmt.Sprintf("/v1/plans/%s", planID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Create creates a new plan
func (p *Plans) Create(ctx context.Context, data CreatePlanRequest) (*Plan, error) {
	var resp DataWrapper[Plan]
	if err := p.http.Post(ctx, "/v1/plans", data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates an existing plan
func (p *Plans) Update(ctx context.Context, planID string, data UpdatePlanRequest) (*Plan, error) {
	var resp DataWrapper[Plan]
	if err := p.http.Post(ctx, fmt.Sprintf("/v1/plans/%s", planID), data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Publish publishes a plan (makes it available for subscriptions)
func (p *Plans) Publish(ctx context.Context, planID string) (*Plan, error) {
	var resp DataWrapper[Plan]
	if err := p.http.Post(ctx, fmt.Sprintf("/v1/plans/%s/publish", planID), struct{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AddFeature adds a feature to a plan
func (p *Plans) AddFeature(ctx context.Context, data CreatePlanFeatureRequest) (*PlanFeature, error) {
	var resp DataWrapper[PlanFeature]
	if err := p.http.Post(ctx, "/v1/plan-features", data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListFeatures lists features for a plan
func (p *Plans) ListFeatures(ctx context.Context, planID string) ([]PlanFeature, error) {
	var resp DataWrapper[[]PlanFeature]
	if err := p.http.Get(ctx, fmt.Sprintf("/v1/plan-features?plan_id=%s", planID), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetFreePlan returns the first free plan (is_free=true)
func (p *Plans) GetFreePlan(ctx context.Context) (*Plan, error) {
	plans, err := p.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, plan := range plans {
		if plan.IsFree {
			return &plan, nil
		}
	}

	return nil, &APIError{Message: "no free plan found"}
}

// GetProPlan returns the first non-free published plan
func (p *Plans) GetProPlan(ctx context.Context) (*Plan, error) {
	plans, err := p.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, plan := range plans {
		if !plan.IsFree && plan.Status == "published" {
			return &plan, nil
		}
	}

	return nil, &APIError{Message: "no pro plan found"}
}
