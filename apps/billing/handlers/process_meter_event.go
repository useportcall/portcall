package handlers

import (
	"encoding/json"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type ProcessMeterEventPayload struct {
	MeterEventID uint `json:"meter_event_id"`
}

func ProcessMeterEvent(c server.IContext) error {
	var p ProcessMeterEventPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var meterEvent models.MeterEvent
	if err := c.DB().FindForID(p.MeterEventID, &meterEvent); err != nil {
		return err
	}

	var feature models.Feature
	if err := c.DB().FindForID(meterEvent.FeatureID, &feature); err != nil {
		return err
	}

	var entitlement models.Entitlement
	if err := c.DB().FindFirst(&entitlement, "user_id = ? AND feature_public_id = ?", meterEvent.UserID, feature.PublicID); err != nil {
		return err
	}

	if err := c.DB().IncrementCount(&entitlement, "usage", meterEvent.Usage); err != nil {
		return err
	}

	return nil
}
