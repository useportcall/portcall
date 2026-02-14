package plan

import (
	"time"

	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CopyPlanRequest struct {
	TargetAppID string `json:"target_app_id"`
}

func CopyPlan(c *routerx.Context) {
	planID := c.Param("id")

	var req CopyPlanRequest
	if err := c.BindJSON(&req); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if req.TargetAppID == "" {
		c.BadRequest("target_app_id is required")
		return
	}

	// Get the source plan
	var sourcePlan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), planID, &sourcePlan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	// Verify target app exists and belongs to the same account
	var targetApp models.App
	if err := c.DB().FindFirst(&targetApp, "public_id = ?", req.TargetAppID); err != nil {
		c.NotFound("Target app not found")
		return
	}

	// Verify both apps belong to the same account
	var sourceApp models.App
	if err := c.DB().FindForID(c.AppID(), &sourceApp); err != nil {
		c.ServerError("Failed to get source app", err)
		return
	}

	if sourceApp.AccountID != targetApp.AccountID {
		c.BadRequest("Target app does not belong to the same account")
		return
	}

	// Create new plan for target app
	newPlan := models.Plan{
		PublicID:         utils.GenPublicID("plan"),
		AppID:            targetApp.ID,
		Name:             sourcePlan.Name,
		Status:           "draft", // Always create as draft
		TrialPeriodDays:  sourcePlan.TrialPeriodDays,
		Interval:         sourcePlan.Interval,
		IntervalCount:    sourcePlan.IntervalCount,
		Currency:         sourcePlan.Currency,
		InvoiceDueByDays: sourcePlan.InvoiceDueByDays,
		IsFree:           sourcePlan.IsFree,
		DiscountPct:      sourcePlan.DiscountPct,
		DiscountQty:      sourcePlan.DiscountQty,
		PlanGroupID:      nil, // Don't copy plan group reference
	}
	newPlan.CreatedAt = time.Now()
	newPlan.UpdatedAt = time.Now()

	if err := c.DB().Create(&newPlan); err != nil {
		c.ServerError("Failed to create plan", err)
		return
	}

	// Copy plan items
	var planItems []models.PlanItem
	if err := c.DB().List(&planItems, "plan_id = ?", sourcePlan.ID); err != nil {
		c.ServerError("Failed to list plan items", err)
		return
	}

	for _, item := range planItems {
		newItem := models.PlanItem{
			PublicID:          utils.GenPublicID("plan_item"),
			AppID:             targetApp.ID,
			PlanID:            newPlan.ID,
			PricingModel:      item.PricingModel,
			Quantity:          item.Quantity,
			UnitAmount:        item.UnitAmount,
			Maximum:           item.Maximum,
			Minimum:           item.Minimum,
			Tiers:             item.Tiers,
			PublicTitle:       item.PublicTitle,
			PublicDescription: item.PublicDescription,
			PublicUnitLabel:   item.PublicUnitLabel,
		}
		newItem.CreatedAt = time.Now()
		newItem.UpdatedAt = time.Now()

		if err := c.DB().Create(&newItem); err != nil {
			c.ServerError("Failed to create plan item", err)
			return
		}

		// Copy plan features for this plan item
		var planFeatures []models.PlanFeature
		if err := c.DB().List(&planFeatures, "plan_item_id = ?", item.ID); err != nil {
			c.ServerError("Failed to list plan features", err)
			return
		}

		for _, pf := range planFeatures {
			// Get the source feature to check/copy it
			var sourceFeature models.Feature
			if pf.FeatureID != 0 {
				if err := c.DB().FindForID(pf.FeatureID, &sourceFeature); err != nil {
					continue // Skip if feature not found
				}

				// Check if feature already exists in target app by public_id
				var targetFeature models.Feature
				err := c.DB().FindFirst(&targetFeature, "app_id = ? AND public_id = ?", targetApp.ID, sourceFeature.PublicID)

				if err != nil {
					// Feature doesn't exist in target app, create it
					targetFeature = models.Feature{
						PublicID:  sourceFeature.PublicID,
						AppID:     targetApp.ID,
						IsMetered: sourceFeature.IsMetered,
					}
					targetFeature.CreatedAt = time.Now()
					targetFeature.UpdatedAt = time.Now()

					if err := c.DB().Create(&targetFeature); err != nil {
						c.ServerError("Failed to create feature", err)
						return
					}
				}

				// Create plan feature association
				newPlanFeature := models.PlanFeature{
					PublicID:   utils.GenPublicID("plan_feature"),
					AppID:      targetApp.ID,
					PlanID:     newPlan.ID,
					PlanItemID: newItem.ID,
					FeatureID:  targetFeature.ID,
					Interval:   pf.Interval,
					Quota:      pf.Quota,
					Rollover:   pf.Rollover,
				}
				newPlanFeature.CreatedAt = time.Now()
				newPlanFeature.UpdatedAt = time.Now()

				if err := c.DB().Create(&newPlanFeature); err != nil {
					c.ServerError("Failed to create plan feature", err)
					return
				}
			}
		}
	}

	// Return the new plan
	response := new(apix.Plan)
	c.OK(response.Set(&newPlan))
}
