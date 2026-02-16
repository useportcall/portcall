package subscription

import (
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
)

// QueueEnqueuer defines the queue capability needed by the billing flow.
type QueueEnqueuer interface {
	Enqueue(name string, payload any, queue string) error
}

// UpgradeNotifier handles side effects for plan upgrades.
type UpgradeNotifier func(userID, newPlanID uint)

// BillingFlowInput is the input for running create/update subscription billing.
type BillingFlowInput struct {
	DB     dbx.IORM
	Crypto cryptox.ICrypto
	Queue  QueueEnqueuer
	AppID  uint
	UserID uint
	PlanID uint
	// AsyncPayment enqueues pay tasks instead of paying inline.
	AsyncPayment  bool
	NotifyUpgrade UpgradeNotifier
}

// BillingFlowResult reports whether a new subscription was created.
type BillingFlowResult struct {
	CreatedNew bool
}
