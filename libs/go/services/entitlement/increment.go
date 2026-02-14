package entitlement

import "github.com/useportcall/portcall/libs/go/dbx/models"

// IncrementUsage processes a meter event by incrementing the related
// entitlement's usage counter. Per-use entitlements are skipped.
func (s *service) IncrementUsage(input *IncrementUsageInput) (*IncrementUsageResult, error) {
	var meterEvent models.MeterEvent
	if err := s.db.FindForID(input.MeterEventID, &meterEvent); err != nil {
		return nil, err
	}

	var feature models.Feature
	if err := s.db.FindForID(meterEvent.FeatureID, &feature); err != nil {
		return nil, err
	}

	var ent models.Entitlement
	if err := s.db.FindFirst(&ent,
		"user_id = ? AND feature_public_id = ?",
		meterEvent.UserID, feature.PublicID,
	); err != nil {
		return nil, err
	}

	if ent.Interval != "per_use" {
		ent.Usage += meterEvent.Usage
		if ent.Usage < 0 {
			ent.Usage = 0
		}
		if err := s.db.Save(&ent); err != nil {
			return nil, err
		}
	}

	if err := syncMeteredUsage(s.db, &meterEvent, &feature); err != nil {
		return nil, err
	}

	result := &IncrementUsageResult{MeterEvent: &meterEvent, Entitlement: &ent}
	result.Skipped = ent.Interval == "per_use"
	return result, nil
}
