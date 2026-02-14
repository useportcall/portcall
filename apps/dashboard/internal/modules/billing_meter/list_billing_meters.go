package billing_meter

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// ListBillingMeters returns all billing meters for a subscription or user.
func ListBillingMeters(c *routerx.Context) {
	subscriptionID := c.Query("subscription_id")
	userID := c.Query("user_id")

	if subscriptionID == "" && userID == "" {
		c.BadRequest("Either subscription_id or user_id is required")
		return
	}

	var meters []models.BillingMeter
	var err error

	if subscriptionID != "" {
		// Get subscription by public ID
		var subscription models.Subscription
		if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
			c.NotFound("Subscription not found")
			return
		}
		err = c.DB().List(&meters, "subscription_id = ?", subscription.ID)
	} else {
		// Get user by public ID
		var user models.User
		if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
			c.NotFound("User not found")
			return
		}
		err = c.DB().List(&meters, "user_id = ?", user.ID)
	}

	if err != nil {
		c.ServerError("Failed to list billing meters", err)
		return
	}

	// Build response with additional info
	response := make([]apix.BillingMeter, len(meters))
	for i, meter := range meters {
		response[i].Set(&meter)

		// Get subscription public ID
		var subscription models.Subscription
		if err := c.DB().FindForID(meter.SubscriptionID, &subscription); err == nil {
			response[i].SubscriptionID = subscription.PublicID
		}

		// Get feature public ID
		var feature models.Feature
		if err := c.DB().FindForID(meter.FeatureID, &feature); err == nil {
			response[i].FeatureID = feature.PublicID
			response[i].FeatureName = feature.PublicID // Feature uses PublicID as name
		}

		// Get plan item info
		var planItem models.PlanItem
		if err := c.DB().FindForID(meter.PlanItemID, &planItem); err == nil {
			response[i].PlanItemID = planItem.PublicID
			response[i].PlanItemTitle = planItem.PublicTitle
		}

		// Get user public ID
		var user models.User
		if err := c.DB().FindForID(meter.UserID, &user); err == nil {
			response[i].UserID = user.PublicID
		}

		// Calculate projected cost
		response[i].ProjectedCost = response[i].CalculateProjectedCost()
	}

	c.OK(response)
}

// GetBillingMeter returns a specific billing meter by subscription and feature.
func GetBillingMeter(c *routerx.Context) {
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
