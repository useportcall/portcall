package meter_event

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateMeterEventRequest struct {
	UserID    string `json:"user_id"`
	FeatureID string `json:"feature_id"`
	Usage     int64  `json:"usage"`
}

func CreateMeterEvent(c *routerx.Context) {
	var body CreateMeterEventRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request payload")
		return
	}

	if body.UserID == "" || body.FeatureID == "" {
		c.BadRequest("user_id and feature_id are required")
		return
	}
	if body.Usage <= 0 {
		c.BadRequest("usage must be a positive integer")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), body.UserID, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	var entitlement models.Entitlement
	if err := c.DB().FindFirst(&entitlement, "app_id = ? AND user_id = ? AND feature_public_id = ?", c.AppID(), user.ID, body.FeatureID); err != nil {
		c.NotFound("Entitlement not found")
		return
	}

	if entitlement.Quota > 0 && (entitlement.Usage+body.Usage) > entitlement.Quota {
		c.BadRequest("Usage exceeds entitlement quota")
		return
	}

	var feature models.Feature
	if err := c.DB().GetForPublicID(c.AppID(), body.FeatureID, &feature); err != nil {
		c.NotFound("Feature not found")
		return
	}

	event := models.MeterEvent{
		AppID:     c.AppID(),
		UserID:    user.ID,
		FeatureID: feature.ID,
		Usage:     body.Usage,
		Timestamp: time.Now(),
	}
	if err := c.DB().Create(&event); err != nil {
		c.ServerError("Failed to create meter event", err)
		return
	}

	if err := c.Queue().Enqueue("process_meter_event", map[string]any{"meter_event_id": event.ID}, "billing_queue"); err != nil {
		c.ServerError("Failed to enqueue meter event processing", err)
		return
	}

	c.OK(nil)
}
