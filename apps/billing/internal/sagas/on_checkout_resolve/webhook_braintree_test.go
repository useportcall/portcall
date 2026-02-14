package on_checkout_resolve_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_checkout_resolve"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/services/payment"
)

func TestWebhook_BraintreeSettled_EnqueuesResolveSession(t *testing.T) {
	db := saga.NewStubDB()
	runner := saga.NewRunner(db, nil, []saga.Step{on_checkout_resolve.Steps[1]})

	err := runner.Run("process_braintree_webhook_event", payment.BraintreeWebhookInput{
		Kind:               "transaction_settled",
		OrderID:            "portcall_checkout_session_id=btsess_1",
		PaymentMethodToken: "token_1",
	})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !runner.HasTask("resolve_checkout_session") {
		t.Fatalf("expected resolve_checkout_session task, got %v", runner.Executed)
	}
}

func TestWebhook_BraintreeDeclined_EnqueuesDunningTask(t *testing.T) {
	db := saga.NewStubDB()
	runner := saga.NewRunner(db, nil, []saga.Step{on_checkout_resolve.Steps[1]})

	err := runner.Run("process_braintree_webhook_event", payment.BraintreeWebhookInput{
		Kind:          "transaction_settlement_declined",
		OrderID:       "portcall_invoice_id=52",
		FailureCount:  3,
		FailureReason: "declined",
	})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !runner.HasTask("process_stripe_payment_failure") {
		t.Fatalf("expected process_stripe_payment_failure task, got %v", runner.Executed)
	}
}
