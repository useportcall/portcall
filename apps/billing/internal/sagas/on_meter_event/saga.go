package on_meter_event

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
)

// --- Flow: process_meter_event (terminal) ---

var ProcessMeterEvent = saga.Route{Name: "process_meter_event", Queue: saga.BillingQueue}

var Steps = []saga.Step{
	{Route: ProcessMeterEvent, Handler: processHandler},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }

func processHandler(c server.IContext) error {
	var input entitlement.IncrementUsageInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := entitlement.NewService(c.DB())
	_, err := svc.IncrementUsage(&input)
	return err
}
