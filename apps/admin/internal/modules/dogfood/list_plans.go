package dogfood

import (
	"strconv"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type PlanDetailResponse struct {
	ID              uint                  `json:"id"`
	PublicID        string                `json:"public_id"`
	Name            string                `json:"name"`
	Description     string                `json:"description,omitempty"`
	IsFree          bool                  `json:"is_free"`
	Interval        string                `json:"interval"`
	IntervalCount   int                   `json:"interval_count"`
	SubscriberCount int64                 `json:"subscriber_count"`
	Features        []PlanFeatureResponse `json:"features"`
}

type PlanFeatureResponse struct {
	FeatureID string `json:"feature_id"`
	Quota     int64  `json:"quota"`
	Interval  string `json:"interval"`
}

// ListPlans lists all plans for a dogfood app with subscriber counts and features
func ListPlans(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	if appIDStr == "" {
		c.BadRequest("Missing app_id parameter")
		return
	}

	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		c.BadRequest("Invalid app_id parameter")
		return
	}

	// Verify the app belongs to dogfood account
	var app models.App
	if err := c.DB().FindForID(uint(appID), &app); err != nil {
		c.NotFound("App not found")
		return
	}

	var account models.Account
	if err := c.DB().FindForID(app.AccountID, &account); err != nil || account.Email != DogfoodAccountEmail {
		c.Unauthorized("Access denied - not a dogfood app")
		return
	}

	// Get all plans for this app
	var plans []models.Plan
	if err := c.DB().List(&plans, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list plans", err)
		return
	}

	// Get plan features for all plans
	planFeatureMap := make(map[uint][]PlanFeatureResponse)
	if len(plans) > 0 {
		planIDs := make([]uint, len(plans))
		for i, p := range plans {
			planIDs[i] = p.ID
		}

		var planFeatures []models.PlanFeature
		if err := c.DB().List(&planFeatures, "plan_id IN ?", planIDs); err == nil {
			// Get feature public IDs
			featureIDs := make([]uint, 0)
			for _, pf := range planFeatures {
				featureIDs = append(featureIDs, pf.FeatureID)
			}

			featureMap := make(map[uint]string)
			if len(featureIDs) > 0 {
				var features []models.Feature
				if err := c.DB().List(&features, "id IN ?", featureIDs); err == nil {
					for _, f := range features {
						featureMap[f.ID] = f.PublicID
					}
				}
			}

			for _, pf := range planFeatures {
				featurePublicID := featureMap[pf.FeatureID]
				if featurePublicID == "" {
					continue
				}
				planFeatureMap[pf.PlanID] = append(planFeatureMap[pf.PlanID], PlanFeatureResponse{
					FeatureID: featurePublicID,
					Quota:     pf.Quota,
					Interval:  pf.Interval,
				})
			}
		}
	}

	// Count subscribers per plan
	subscriberCounts := make(map[uint]int64)
	for _, p := range plans {
		var count int64
		c.DB().Count(&count, &models.Subscription{}, "plan_id = ? AND status = ?", p.ID, "active")
		subscriberCounts[p.ID] = count
	}

	// Build response
	result := make([]PlanDetailResponse, len(plans))
	for i, p := range plans {
		result[i] = PlanDetailResponse{
			ID:              p.ID,
			PublicID:        p.PublicID,
			Name:            p.Name,
			IsFree:          p.IsFree,
			Interval:        p.Interval,
			IntervalCount:   p.IntervalCount,
			SubscriberCount: subscriberCounts[p.ID],
			Features:        planFeatureMap[p.ID],
		}
		if result[i].Features == nil {
			result[i].Features = []PlanFeatureResponse{}
		}
	}

	c.OK(result)
}
