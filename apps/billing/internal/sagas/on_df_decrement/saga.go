package on_df_decrement

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/df"
)

// --- Flow: df_decrement (terminal) ---

var Decrement = saga.Route{Name: "df_decrement", Queue: saga.BillingQueue}

var Steps = []saga.Step{
	{Route: Decrement, Handler: decrementHandler},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }

func decrementHandler(c server.IContext) error {
	var input df.DecrementInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := df.NewService()
	return svc.Decrement(c, &input)
}
