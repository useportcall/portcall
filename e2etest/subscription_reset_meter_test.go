package e2etest

import (
	"testing"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
)

func processMeterEvent(t *testing.T, h *Harness, userID, featureID string, usage int64) {
	t.Helper()
	h.RecordMeterEventViaAPI(userID, featureID, int(usage)).MustStatus(t, 200, "record meter event")
	var user models.User
	must(t, h.DB.GetForPublicID(h.AppID, userID, &user))
	var feature models.Feature
	must(t, h.DB.GetForPublicID(h.AppID, featureID, &feature))
	var event models.MeterEvent
	must(t, h.DB.FindFirst(&event, "app_id = ? AND user_id = ? AND feature_id = ?", h.AppID, user.ID, feature.ID))
	r := billing.NewRunner(h.DB, h.Crypto, billing.MeterEventSteps())
	must(t, r.Run("process_meter_event", entitlement.IncrementUsageInput{MeterEventID: event.ID}))
}
