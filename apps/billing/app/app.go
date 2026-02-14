// Package app exposes billing saga types for cross-app e2e testing.
package app

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_checkout_resolve"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_entitlement_reset"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_meter_event"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_reset"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_update"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
)

// Step re-exports saga.Step for external consumers.
type Step = saga.Step

// Runner re-exports saga.Runner for external consumers.
type Runner = saga.Runner

// NewRunner creates a saga runner with the given DB, crypto, and step sets.
var NewRunner = saga.NewRunner

// AllSteps returns every billing saga step collection for a full chain test.
func AllSteps() [][]Step {
	return [][]Step{
		on_checkout_resolve.Steps,
		on_subscription_create.Steps,
		on_subscription_update.Steps,
		on_invoice_create.Steps,
		on_payment.Steps,
		on_subscription_reset.Steps,
		on_entitlement_reset.Steps,
	}
}

// AllStepsWithMeter returns all billing steps plus meter event processing.
func AllStepsWithMeter() [][]Step {
	return [][]Step{
		on_checkout_resolve.Steps,
		on_subscription_create.Steps,
		on_subscription_update.Steps,
		on_invoice_create.Steps,
		on_payment.Steps,
		on_subscription_reset.Steps,
		on_entitlement_reset.Steps,
		on_meter_event.Steps,
	}
}

// ResetSteps returns subscription reset saga steps.
func ResetSteps() []Step {
	return on_subscription_reset.Steps
}

// MeterEventSteps returns only the meter event saga steps.
func MeterEventSteps() []Step {
	return on_meter_event.Steps
}

// NewFullRunner creates a saga runner pre-loaded with all billing steps.
func NewFullRunner(db dbx.IORM, crypto cryptox.ICrypto) *Runner {
	steps := AllSteps()
	return saga.NewRunner(db, crypto, steps...)
}
