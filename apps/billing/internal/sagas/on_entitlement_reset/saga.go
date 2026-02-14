package on_entitlement_reset

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
)

// --- Flow: reset_all_entitlements (single bulk operation) ---

var ResetAll = saga.Route{Name: "reset_all_entitlements", Queue: saga.BillingQueue}

var Steps = []saga.Step{
	{Route: ResetAll, Handler: resetAllHandler, Emits: nil},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }

func resetAllHandler(c server.IContext) error {
	var input entitlement.ResetAllInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := entitlement.NewService(c.DB())
	_, err := svc.ResetAll(&input)
	return err
}
