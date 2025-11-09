package feature

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Feature struct {
	ID        string    `json:"id"`
	IsMetered bool      `json:"is_metered"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Feature) Set(feature *models.Feature) *Feature {
	f.ID = feature.PublicID
	f.IsMetered = feature.IsMetered
	f.CreatedAt = feature.CreatedAt
	f.UpdatedAt = feature.UpdatedAt
	return f
}

type CreateFeatureRequest struct {
	FeatureID     string `json:"feature_id"`
	IsMetered     bool   `json:"is_metered"`
	PlanID        string `json:"plan_id,omitempty"`
	PlanFeatureID string `json:"plan_feature_id,omitempty"`
}
