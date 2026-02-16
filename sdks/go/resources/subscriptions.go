package resources

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Subscription represents a user's subscription
type Subscription struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	PlanID             *string    `json:"plan_id,omitempty"`
	Plan               *Plan      `json:"plan,omitempty"`
	ScheduledPlanID    *string    `json:"scheduled_plan_id,omitempty"`
	ScheduledPlan      *Plan      `json:"scheduled_plan,omitempty"`
	Status             string     `json:"status"`
	IsFree             bool       `json:"is_free"`
	CurrentPeriodStart *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`
	NextResetAt        *time.Time `json:"next_reset_at,omitempty"`
	CancelAtPeriodEnd  *bool      `json:"cancel_at_period_end,omitempty"`
	CanceledAt         *time.Time `json:"canceled_at,omitempty"`
	TrialStart         *time.Time `json:"trial_start,omitempty"`
	TrialEnd           *time.Time `json:"trial_end,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ListSubscriptionsParams are the parameters for listing subscriptions
type ListSubscriptionsParams struct {
	UserID string
	Status string
	Limit  int
}

// UpdateSubscriptionRequest is the request to update a subscription
type UpdateSubscriptionRequest struct {
	PlanID            *string                `json:"plan_id,omitempty"`
	CancelAtPeriodEnd *bool                  `json:"cancel_at_period_end,omitempty"`
	ApplyAtNextReset  *bool                  `json:"apply_at_next_reset,omitempty"` // If true, schedule plan change for next reset
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// CancelSubscriptionRequest is the request to cancel a subscription
type CancelSubscriptionRequest struct {
	Immediately bool `json:"immediately,omitempty"`
}

// CreateSubscriptionRequest is the request to create a subscription
type CreateSubscriptionRequest struct {
	UserID string `json:"user_id"`
	PlanID string `json:"plan_id"`
}

// Subscriptions provides access to subscription-related API operations
type Subscriptions struct {
	http *HTTPClient
}

// NewSubscriptions creates a new Subscriptions resource
func NewSubscriptions(http *HTTPClient) *Subscriptions {
	return &Subscriptions{http: http}
}

// List returns all subscriptions matching the given parameters
func (s *Subscriptions) List(ctx context.Context, params *ListSubscriptionsParams) ([]Subscription, error) {
	query := url.Values{}
	if params != nil {
		if params.UserID != "" {
			query.Set("user_id", params.UserID)
		}
		if params.Status != "" {
			query.Set("status", params.Status)
		}
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
	}

	path := "/v1/subscriptions"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var resp DataWrapper[[]Subscription]
	if err := s.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a subscription by ID
func (s *Subscriptions) Get(ctx context.Context, subscriptionID string) (*Subscription, error) {
	var resp DataWrapper[Subscription]
	if err := s.http.Get(ctx, fmt.Sprintf("/v1/subscriptions/%s", subscriptionID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Create creates a new subscription
func (s *Subscriptions) Create(ctx context.Context, data CreateSubscriptionRequest) (*Subscription, error) {
	var resp DataWrapper[Subscription]
	if err := s.http.Post(ctx, "/v1/subscriptions", data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates an existing subscription
func (s *Subscriptions) Update(ctx context.Context, subscriptionID string, data UpdateSubscriptionRequest) (*Subscription, error) {
	var resp DataWrapper[Subscription]
	if err := s.http.Post(ctx, fmt.Sprintf("/v1/subscriptions/%s", subscriptionID), data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Cancel cancels a subscription
func (s *Subscriptions) Cancel(ctx context.Context, subscriptionID string, data *CancelSubscriptionRequest) (*Subscription, error) {
	var resp DataWrapper[Subscription]
	if data == nil {
		data = &CancelSubscriptionRequest{}
	}
	if err := s.http.Post(ctx, fmt.Sprintf("/v1/subscriptions/%s/cancel", subscriptionID), data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetActive returns the active subscription for a user, if any
func (s *Subscriptions) GetActive(ctx context.Context, userID string) (*Subscription, error) {
	subs, err := s.List(ctx, &ListSubscriptionsParams{
		UserID: userID,
		Status: "active",
		Limit:  1,
	})
	if err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return nil, nil
	}

	return &subs[0], nil
}
