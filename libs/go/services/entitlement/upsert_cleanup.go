package entitlement

import "github.com/useportcall/portcall/libs/go/dbx/models"

func removeMissingEntitlements(s *service, userID uint, planFeatureIDs []uint) error {
	allowed, err := allowedFeatureSet(s, planFeatureIDs)
	if err != nil {
		return err
	}

	var entitlementIDs []uint
	if err := s.db.ListIDs("entitlements", &entitlementIDs, "user_id = ?", userID); err != nil {
		return err
	}
	for _, id := range entitlementIDs {
		var ent models.Entitlement
		if err := s.db.FindForID(id, &ent); err != nil {
			continue
		}
		if allowed[ent.FeaturePublicID] {
			continue
		}
		if err := s.db.Delete(&models.Entitlement{}, "id = ?", ent.ID); err != nil {
			return err
		}
	}
	return nil
}

func allowedFeatureSet(s *service, planFeatureIDs []uint) (map[string]bool, error) {
	allowed := map[string]bool{}
	for _, id := range planFeatureIDs {
		var pf models.PlanFeature
		if err := s.db.FindForID(id, &pf); err != nil {
			return nil, err
		}
		var feature models.Feature
		if err := s.db.FindForID(pf.FeatureID, &feature); err != nil {
			return nil, err
		}
		allowed[feature.PublicID] = true
	}
	return allowed, nil
}
