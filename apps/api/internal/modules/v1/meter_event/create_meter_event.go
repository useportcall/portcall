package meter_event

import (
	"log"
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
		log.Println("Invalid request payload:", err)
		c.ServerError("Invalid request payload", err)
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), body.UserID, &user); err != nil {
		log.Println("User not found:", err)
		c.ServerError("User not found", err)
		return
	}

	var entitlement models.Entitlement
	if err := c.DB().FindFirst(&entitlement, "app_id = ? AND feature_public_id = ?", c.AppID(), body.FeatureID); err != nil {
		c.ServerError("Entitlement not found", err)
		return
	}

	log.Println("Current entitlement usage:", entitlement.Usage, "Quota:", entitlement.Quota, "Requested usage increment:", body.Usage)

	if entitlement.Quota > 0 && (entitlement.Usage+body.Usage) > entitlement.Quota && body.Usage > 0 {
		c.BadRequest("Usage exceeds entitlement quota")
		return
	}

	var feature models.Feature
	if err := c.DB().GetForPublicID(c.AppID(), body.FeatureID, &feature); err != nil {
		log.Println("Feature not found:", err)
		c.ServerError("Feature not found", err)
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
		log.Println("Error creating meter event:", err)
		c.ServerError("Failed to create meter event", err)
		return
	}

	if err := c.Queue().Enqueue("process_meter_event", map[string]any{"meter_event_id": event.ID}, "billing_queue"); err != nil {
		log.Println("Error enqueuing meter event processing:", err)
		c.ServerError("Failed to enqueue meter event processing", err)
		return
	}

	c.OK(nil)
}
