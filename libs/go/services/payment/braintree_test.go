package payment_test

import (
	"testing"

	pay "github.com/useportcall/portcall/libs/go/services/payment"
)

func TestProcessBraintreeWebhook_ResolveCheckoutSession(t *testing.T) {
	svc := pay.NewService(nil, nil)
	result, err := svc.ProcessBraintreeWebhook(&pay.BraintreeWebhookInput{
		Kind:               "transaction_settled",
		OrderID:            "portcall_checkout_session_id=btsess_123",
		PaymentMethodToken: "token_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Handled || result.Action != "resolve_checkout_session" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestProcessBraintreeWebhook_FailurePayload(t *testing.T) {
	svc := pay.NewService(nil, nil)
	result, err := svc.ProcessBraintreeWebhook(&pay.BraintreeWebhookInput{
		Kind:          "transaction_settlement_declined",
		OrderID:       "portcall_invoice_id=42",
		FailureCount:  2,
		FailureReason: "declined",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Failure == nil || result.Failure.InvoiceID != 42 || result.Failure.Attempt != 2 {
		t.Fatalf("unexpected failure payload: %+v", result.Failure)
	}
}

func TestProcessBraintreeWebhook_IgnoresUnknown(t *testing.T) {
	svc := pay.NewService(nil, nil)
	result, err := svc.ProcessBraintreeWebhook(&pay.BraintreeWebhookInput{Kind: "dispute_opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Handled {
		t.Fatalf("expected unhandled result, got %+v", result)
	}
}
