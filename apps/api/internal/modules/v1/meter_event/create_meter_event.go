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
		c.ServerError("Invalid request payload")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), body.UserID, &user); err != nil {
		log.Println("User not found:", err)
		c.ServerError("User not found")
		return
	}

	var feature models.Feature
	if err := c.DB().GetForPublicID(c.AppID(), body.FeatureID, &feature); err != nil {
		log.Println("Feature not found:", err)
		c.ServerError("Feature not found")
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
		c.ServerError("Failed to create meter event")
		return
	}

	if err := c.Queue().Enqueue("process_meter_event", map[string]any{"meter_event_id": event.ID}, "billing_queue"); err != nil {
		log.Println("Error enqueuing meter event processing:", err)
		c.ServerError("Failed to enqueue meter event processing")
		return
	}

	c.OK(nil)
}
