package payment

import "testing"

func TestParseOrderValue_Keyed(t *testing.T) {
	got := parseOrderMetadataValue(
		"portcall_checkout_session_id=btsess_1",
		"portcall_checkout_session_id",
	)
	if got != "btsess_1" {
		t.Fatalf("got %q, want btsess_1", got)
	}
}

func TestParseOrderValue_Prefix(t *testing.T) {
	got := parseOrderMetadataValue(
		"portcall_invoice_99",
		"portcall_invoice_id",
	)
	if got != "99" {
		t.Fatalf("got %q, want 99", got)
	}
}

func TestParseOrderValue_Empty(t *testing.T) {
	got := parseOrderMetadataValue("", "portcall_invoice_id")
	if got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}

func TestParseOrderUint_Valid(t *testing.T) {
	got := parseOrderMetadataUint(
		"portcall_invoice_id=42",
		"portcall_invoice_id",
	)
	if got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
}

func TestParseOrderUint_Invalid(t *testing.T) {
	got := parseOrderMetadataUint(
		"portcall_invoice_id=abc",
		"portcall_invoice_id",
	)
	if got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
}

func TestParseOrderValue_Pipe(t *testing.T) {
	got := parseOrderMetadataValue(
		"portcall_invoice_id=10|extra=x",
		"portcall_invoice_id",
	)
	if got != "10" {
		t.Fatalf("got %q, want 10", got)
	}
}
