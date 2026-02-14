package payment_test

import (
	"encoding/json"
	"testing"

	"github.com/stripe/stripe-go"
	pay "github.com/useportcall/portcall/libs/go/services/payment"
)

func TestProcessStripeWebhook_SetupIntentSucceeded(t *testing.T) {
	raw := json.RawMessage(`{"id":"seti_123","payment_method":"pm_123"}`)
	svc := pay.NewService(nil, nil)
	result, err := svc.ProcessStripeWebhook(&pay.StripeWebhookInput{
		Event: stripe.Event{Type: "setup_intent.succeeded", Data: &stripe.EventData{Raw: raw}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Handled {
		t.Fatal("expected handled=true")
	}
	if result.Action != "resolve_checkout_session" {
		t.Fatalf("expected resolve_checkout_session, got %s", result.Action)
	}
}

func TestProcessStripeWebhook_UnhandledType(t *testing.T) {
	svc := pay.NewService(nil, nil)
	result, err := svc.ProcessStripeWebhook(&pay.StripeWebhookInput{
		Event: stripe.Event{Type: "payment_method.attached"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Handled {
		t.Fatal("expected handled=false")
	}
}

func TestProcessStripeWebhook_PaymentIntentFailed_WithInvoiceMetadata(t *testing.T) {
	raw := json.RawMessage(`{
		"id":"pi_123",
		"metadata":{"portcall_invoice_id":"77"},
		"last_payment_error":{"decline_code":"lost_card","message":"Card lost"}
	}`)
	svc := pay.NewService(nil, nil)
	result, err := svc.ProcessStripeWebhook(&pay.StripeWebhookInput{
		Event: stripe.Event{
			Type: "payment_intent.payment_failed",
			Data: &stripe.EventData{Raw: raw},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "process_stripe_payment_failure" {
		t.Fatalf("expected failure action, got %s", result.Action)
	}
	if result.Failure == nil || result.Failure.InvoiceID != 77 {
		t.Fatalf("unexpected failure payload: %+v", result.Failure)
	}
	if !result.Failure.NoRetry {
		t.Fatal("expected hard-decline no_retry=true")
	}
}

func TestProcessStripeWebhook_PaymentIntentCanceled_WithoutLastError(t *testing.T) {
	raw := json.RawMessage(`{
		"id":"pi_456",
		"metadata":{"portcall_invoice_id":"88"}
	}`)
	svc := pay.NewService(nil, nil)
	result, err := svc.ProcessStripeWebhook(&pay.StripeWebhookInput{
		Event: stripe.Event{
			Type: "payment_intent.canceled",
			Data: &stripe.EventData{Raw: raw},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "process_stripe_payment_failure" {
		t.Fatalf("expected failure action, got %s", result.Action)
	}
	if result.Failure == nil || result.Failure.InvoiceID != 88 {
		t.Fatalf("unexpected failure payload: %+v", result.Failure)
	}
	if result.Failure.NoRetry {
		t.Fatal("expected no_retry=false when decline code is unavailable")
	}
}
