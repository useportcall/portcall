package resources

import (
	"context"
	"fmt"
	"time"
)

// Feature represents a feature definition
type Feature struct {
	ID          string    `json:"id"`
	Name        string    `json:"name,omitempty"`
	Key         string    `json:"key,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description *string   `json:"description,omitempty"`
	IsMetered   bool      `json:"is_metered"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateFeatureRequest is the request to create a feature
type CreateFeatureRequest struct {
	FeatureID string `json:"feature_id"`
	IsMetered bool   `json:"is_metered,omitempty"`
}

// ListFeaturesParams are the parameters for listing features
type ListFeaturesParams struct {
	IsMetered *bool
}

// Features provides access to feature-related API operations
type Features struct {
	http *HTTPClient
}

// NewFeatures creates a new Features resource
func NewFeatures(http *HTTPClient) *Features {
	return &Features{http: http}
}

// List returns all features
func (f *Features) List(ctx context.Context, params *ListFeaturesParams) ([]Feature, error) {
	url := "/v1/features"
	if params != nil && params.IsMetered != nil {
		url = fmt.Sprintf("/v1/features?is_metered=%v", *params.IsMetered)
	}

	var resp DataWrapper[[]Feature]
	if err := f.http.Get(ctx, url, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ListStandard returns all non-metered (boolean) features
func (f *Features) ListStandard(ctx context.Context) ([]Feature, error) {
	isMetered := false
	return f.List(ctx, &ListFeaturesParams{IsMetered: &isMetered})
}

// ListMetered returns all metered features
func (f *Features) ListMetered(ctx context.Context) ([]Feature, error) {
	isMetered := true
	return f.List(ctx, &ListFeaturesParams{IsMetered: &isMetered})
}

// Create creates a new feature
func (f *Features) Create(ctx context.Context, data CreateFeatureRequest) (*Feature, error) {
	var resp DataWrapper[Feature]
	if err := f.http.Post(ctx, "/v1/features", data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
