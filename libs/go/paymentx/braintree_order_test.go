package paymentx

import "testing"

func TestBraintreeOrderID_InvoiceID(t *testing.T) {
	meta := map[string]string{"portcall_invoice_id": "42"}
	id := braintreeOrderID(meta)
	if id != "portcall_invoice_id=42" {
		t.Fatalf("got %q, want portcall_invoice_id=42", id)
	}
}

func TestBraintreeOrderID_SessionID(t *testing.T) {
	meta := map[string]string{"portcall_checkout_session_id": "cs_abc"}
	id := braintreeOrderID(meta)
	if id != "portcall_checkout_session_id=cs_abc" {
		t.Fatalf("got %q, want portcall_checkout_session_id=cs_abc", id)
	}
}

func TestBraintreeOrderID_Empty(t *testing.T) {
	id := braintreeOrderID(map[string]string{})
	if id != "" {
		t.Fatalf("got %q, want empty", id)
	}
}
