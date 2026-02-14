package on_checkout_resolve_test

import (
	"encoding/json"
	"testing"

	"github.com/stripe/stripe-go"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_checkout_resolve"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/services/payment"
)

func TestWebhook_SetupIntent_EnqueuesResolveSession(t *testing.T) {
	raw := json.RawMessage(`{"id":"seti_webhook_1","payment_method":{"id":"pm_webhook_1"}}`)
	db := saga.NewStubDB()
	runner := saga.NewRunner(db, nil, []saga.Step{on_checkout_resolve.Steps[0]})

	err := runner.Run("process_stripe_webhook_event", payment.StripeWebhookInput{
		Event: stripe.Event{
			Type: "setup_intent.succeeded",
			Data: &stripe.EventData{Raw: raw},
		},
	})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !runner.HasTask("resolve_checkout_session") {
		t.Fatalf("expected resolve_checkout_session task, got %v", runner.Executed)
	}
}

func TestWebhook_PaymentIntentFailed_EnqueuesDunningTask(t *testing.T) {
	raw := json.RawMessage(`{
		"id":"pi_webhook_fail_1",
		"metadata":{"portcall_invoice_id":"42"},
		"last_payment_error":{"decline_code":"stolen_card","message":"Card reported stolen"}
	}`)
	db := saga.NewStubDB()
	runner := saga.NewRunner(db, nil, []saga.Step{on_checkout_resolve.Steps[0]})

	err := runner.Run("process_stripe_webhook_event", payment.StripeWebhookInput{
		Event: stripe.Event{
			Type: "payment_intent.payment_failed",
			Data: &stripe.EventData{Raw: raw},
		},
	})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !runner.HasTask("process_stripe_payment_failure") {
		t.Fatalf("expected process_stripe_payment_failure task, got %v", runner.Executed)
	}

	var payload payment.StripeFailurePayload
	if err := runner.TaskPayload("process_stripe_payment_failure", &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	if payload.InvoiceID != 42 {
		t.Fatalf("expected invoice_id=42, got %d", payload.InvoiceID)
	}
	if !payload.NoRetry {
		t.Fatal("expected no_retry=true for hard decline")
	}
}

func TestWebhook_InvoicePaymentFailed_PassesAttemptCount(t *testing.T) {
	raw := json.RawMessage(`{
		"id":"in_123",
		"attempt_count":2,
		"metadata":{"portcall_invoice_id":"42"}
	}`)
	db := saga.NewStubDB()
	runner := saga.NewRunner(db, nil, []saga.Step{on_checkout_resolve.Steps[0]})

	err := runner.Run("process_stripe_webhook_event", payment.StripeWebhookInput{
		Event: stripe.Event{
			Type: "invoice.payment_failed",
			Data: &stripe.EventData{Raw: raw},
		},
	})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	var payload payment.StripeFailurePayload
	if err := runner.TaskPayload("process_stripe_payment_failure", &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	if payload.Attempt != 2 {
		t.Fatalf("expected attempt=2, got %d", payload.Attempt)
	}
}
