package billing_meter

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateBillingMeterRequest struct {
	Usage     *int64 `json:"usage"`     // Set absolute usage value
	Increment *int64 `json:"increment"` // Add to current usage
	Decrement *int64 `json:"decrement"` // Subtract from current usage
}

// UpdateBillingMeter updates the usage for a billing meter.
// Supports setting absolute value, incrementing, or decrementing.
func UpdateBillingMeter(c *routerx.Context) {
	subscriptionID := c.Param("subscription_id")
	featureID := c.Param("feature_id")

	var body UpdateBillingMeterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var subscription models.Subscription
	if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
		c.NotFound("Subscription not found")
		return
	}

	var feature models.Feature
	if err := c.DB().GetForPublicID(c.AppID(), featureID, &feature); err != nil {
		c.NotFound("Feature not found")
		return
	}

	var meter models.BillingMeter
	if err := c.DB().FindFirst(&meter, "subscription_id = ? AND feature_id = ?", subscription.ID, feature.ID); err != nil {
		c.NotFound("Billing meter not found")
		return
	}

	// Apply the update
	if body.Usage != nil {
		meter.Usage = *body.Usage
	} else if body.Increment != nil {
		meter.Usage += *body.Increment
	} else if body.Decrement != nil {
		meter.Usage -= *body.Decrement
		if meter.Usage < 0 {
			meter.Usage = 0
		}
	}

	if err := c.DB().Save(&meter); err != nil {
		c.ServerError("Failed to update billing meter", err)
		return
	}

	response := new(apix.BillingMeter).Set(&meter)
	response.SubscriptionID = subscription.PublicID
	response.FeatureID = feature.PublicID
	response.FeatureName = feature.PublicID

	var planItem models.PlanItem
	if err := c.DB().FindForID(meter.PlanItemID, &planItem); err == nil {
		response.PlanItemID = planItem.PublicID
		response.PlanItemTitle = planItem.PublicTitle
	}

	var user models.User
	if err := c.DB().FindForID(meter.UserID, &user); err == nil {
		response.UserID = user.PublicID
	}

	response.ProjectedCost = response.CalculateProjectedCost()

	c.OK(response)
}

// ResetBillingMeter resets a billing meter to zero.
func ResetBillingMeter(c *routerx.Context) {
	subscriptionID := c.Param("subscription_id")
	featureID := c.Param("feature_id")

	var subscription models.Subscription
	if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
		c.NotFound("Subscription not found")
		return
	}

	var feature models.Feature
	if err := c.DB().GetForPublicID(c.AppID(), featureID, &feature); err != nil {
		c.NotFound("Feature not found")
		return
	}

	var meter models.BillingMeter
	if err := c.DB().FindFirst(&meter, "subscription_id = ? AND feature_id = ?", subscription.ID, feature.ID); err != nil {
		c.NotFound("Billing meter not found")
		return
	}

	meter.Usage = 0

	if err := c.DB().Save(&meter); err != nil {
		c.ServerError("Failed to reset billing meter", err)
		return
	}

	response := new(apix.BillingMeter).Set(&meter)
	response.SubscriptionID = subscription.PublicID
	response.FeatureID = feature.PublicID
	response.ProjectedCost = 0

	c.OK(response)
}
