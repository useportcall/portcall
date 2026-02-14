package entitlement

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// Upsert processes a single entitlement upsert. Each call handles exactly
// one entitlement (single mutation), then returns info for continuing iteration.
func (s *service) Upsert(input *UpsertInput) (*UpsertResult, error) {
	if input.Index >= len(input.Values) {
		return &UpsertResult{HasMore: false}, nil
	}

	planFeatureID := input.Values[input.Index]

	var planFeature models.PlanFeature
	if err := s.db.FindForID(planFeatureID, &planFeature); err != nil {
		return nil, err
	}

	var feature models.Feature
	if err := s.db.FindForID(planFeature.FeatureID, &feature); err != nil {
		return nil, err
	}

	ent := buildEntitlement(input.UserID, &feature, &planFeature)
	if err := s.upsertEntitlement(input.UserID, feature.PublicID, ent); err != nil {
		return nil, err
	}

	return &UpsertResult{
		UserID:         input.UserID,
		PlanFeatureIDs: input.Values,
		NextIndex:      input.Index + 1,
		HasMore:        input.Index+1 < len(input.Values),
	}, nil
}

// upsertEntitlement creates or updates an entitlement for a user/feature pair.
func (s *service) upsertEntitlement(userID uint, featurePublicID string, ent *models.Entitlement) error {
	var existing models.Entitlement
	err := s.db.FindFirst(&existing, "user_id = ? AND feature_public_id = ?", userID, featurePublicID)
	if err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			return err
		}
		return s.db.Create(ent)
	}

	existing.Quota = ent.Quota
	existing.Interval = ent.Interval
	existing.IsMetered = ent.IsMetered
	existing.Usage = ent.Usage
	existing.LastResetAt = ent.LastResetAt
	existing.NextResetAt = ent.NextResetAt
	return s.db.Save(&existing)
}

// buildEntitlement creates a new entitlement from a plan feature.
func buildEntitlement(userID uint, feature *models.Feature, pf *models.PlanFeature) *models.Entitlement {
	now := time.Now()
	return &models.Entitlement{
		AppID:           pf.AppID,
		UserID:          userID,
		FeaturePublicID: feature.PublicID,
		Interval:        pf.Interval,
		Quota:           int64(pf.Quota),
		Usage:           0,
		IsMetered:       feature.IsMetered,
		LastResetAt:     &now,
		NextResetAt:     nextResetAt(pf.Interval, now),
		AnchorAt:        &now,
	}
}
