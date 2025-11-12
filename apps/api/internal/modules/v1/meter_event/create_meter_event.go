package meter_event

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateMeterEventRequest struct {
	UserID    uint  `json:"user_id"`
	FeatureID uint  `json:"feature_id"`
	Usage     int64 `json:"usage"`
}

func CreateMeterEvent(c *routerx.Context) {
	var body CreateMeterEventRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.ServerError("Invalid request payload")
		return
	}

	event := models.MeterEvent{
		AppID:     c.AppID(),
		UserID:    body.UserID,
		FeatureID: body.FeatureID,
		Usage:     body.Usage,
		Timestamp: time.Now(),
	}
	if err := c.DB().Create(&event); err != nil {
		c.ServerError("Failed to create meter event")
		return
	}

	if err := c.Queue().Enqueue("process_meter_event", map[string]any{"meter_event_id": event.ID}, "billing_queue"); err != nil {
		c.ServerError("Failed to enqueue meter event processing")
		return
	}

	c.OK(nil)
}
