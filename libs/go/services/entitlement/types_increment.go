package entitlement

import "github.com/useportcall/portcall/libs/go/dbx/models"

// IncrementUsageInput is the input for incrementing entitlement usage.
type IncrementUsageInput struct {
	MeterEventID uint `json:"meter_event_id"`
}

// IncrementUsageResult is the result of incrementing entitlement usage.
type IncrementUsageResult struct {
	MeterEvent  *models.MeterEvent
	Entitlement *models.Entitlement
	Skipped     bool
}
