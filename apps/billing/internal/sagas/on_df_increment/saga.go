package on_df_increment

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/df"
)

// --- Flow: df_increment (terminal) ---

var Increment = saga.Route{Name: "df_increment", Queue: saga.BillingQueue}

var Steps = []saga.Step{
	{Route: Increment, Handler: incrementHandler},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }

func incrementHandler(c server.IContext) error {
	var input df.IncrementInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := df.NewService()
	return svc.Increment(c, &input)
}
