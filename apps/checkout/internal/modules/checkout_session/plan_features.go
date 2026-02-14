package checkout_session

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func loadPlanFeatures(db dbx.IORM, planID uint, planItems []models.PlanItem, plan *apix.Plan) error {
	var planFeatures []models.PlanFeature
	if err := db.List(&planFeatures, "plan_id = ?", planID); err != nil {
		return err
	}

	if len(planFeatures) == 0 {
		return nil
	}

	// Batch-load all features in a single query to avoid N+1.
	featureIDs := make([]uint, len(planFeatures))
	for i, pf := range planFeatures {
		featureIDs[i] = pf.FeatureID
	}
	var features []models.Feature
	if err := db.List(&features, "id IN ?", featureIDs); err != nil {
		return err
	}
	featureMap := make(map[uint]models.Feature, len(features))
	for _, f := range features {
		featureMap[f.ID] = f
	}

	for _, pf := range planFeatures {
		feature, ok := featureMap[pf.FeatureID]
		if !ok {
			continue
		}

		if feature.IsMetered {
			res := apix.PlanFeature{}
			res.Set(&pf)
			res.Feature = apix.Feature{ID: feature.PublicID, IsMetered: feature.IsMetered}
			for _, item := range planItems {
				if item.ID == pf.PlanItemID {
					planItem := apix.PlanItem{}
					planItem.Set(&item)
					res.PlanItem = &planItem
					break
				}
			}
			plan.MeteredFeatures = append(plan.MeteredFeatures, res)
		} else {
			res := apix.Feature{ID: feature.PublicID, IsMetered: feature.IsMetered}
			plan.Features = append(plan.Features, res)
		}
	}
	return nil
}
