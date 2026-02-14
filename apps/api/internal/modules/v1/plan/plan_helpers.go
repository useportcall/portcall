package plan

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// loadPlanItems fetches plan items from the database and returns them as
// API response objects. Returns nil and writes an error response on failure.
func loadPlanItems(c *routerx.Context, planID uint) ([]apix.PlanItem, bool) {
	var items []models.PlanItem
	if err := c.DB().List(&items, "plan_id = ?", planID); err != nil {
		c.ServerError("Failed to list plan items", err)
		return nil, false
	}

	result := make([]apix.PlanItem, len(items))
	for i, item := range items {
		result[i].Set(&item)
	}
	return result, true
}

// loadPlanFeatures fetches plan features from the database and returns them
// as API response objects. Returns nil and writes an error response on failure.
func loadPlanFeatures(c *routerx.Context, planID uint) ([]any, bool) {
	var planFeatures []models.PlanFeature
	if err := c.DB().List(&planFeatures, "plan_id = ?", planID); err != nil {
		c.ServerError("Failed to list plan features", err)
		return nil, false
	}

	result := make([]any, len(planFeatures))
	for i, pf := range planFeatures {
		planFeature := new(apix.PlanFeature).Set(&pf)
		var feature models.Feature
		if err := c.DB().FindForID(pf.FeatureID, &feature); err == nil {
			planFeature.Feature = feature.PublicID
		}
		result[i] = planFeature
	}
	return result, true
}
