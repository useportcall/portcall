package on_entitlement_upsert

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
)

// --- Flow: start_entitlement_upsert → upsert_entitlement (iterative loop) ---

var (
	StartUpsert = saga.Route{Name: "start_entitlement_upsert", Queue: saga.BillingQueue}
	Upsert      = saga.Route{Name: "upsert_entitlement", Queue: saga.BillingQueue}
)

var Steps = []saga.Step{
	{Route: StartUpsert, Handler: startUpsertHandler, Emits: []saga.Route{Upsert}},
	{Route: Upsert, Handler: upsertHandler, Emits: []saga.Route{Upsert}},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }

func startUpsertHandler(c server.IContext) error {
	var input entitlement.StartUpsertInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := entitlement.NewService(c.DB())
	result, err := svc.StartUpsert(&input)
	if err != nil || !result.HasMore {
		return err
	}
	return Upsert.Enqueue(c.Queue(), entitlement.UpsertInput{
		UserID: result.UserID,
		Index:  0,
		Values: result.PlanFeatureIDs,
	})
}

func upsertHandler(c server.IContext) error {
	var input entitlement.UpsertInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := entitlement.NewService(c.DB())
	result, err := svc.Upsert(&input)
	if err != nil || !result.HasMore {
		return err
	}
	return Upsert.Enqueue(c.Queue(), entitlement.UpsertInput{
		UserID: result.UserID,
		Index:  result.NextIndex,
		Values: result.PlanFeatureIDs,
	})
}
