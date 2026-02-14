package payment_link

import "testing"

func TestDeriveUserName(t *testing.T) {
	if got := deriveUserName("", "jane.doe@example.com"); got != "jane.doe" {
		t.Fatalf("expected email local-part fallback, got %q", got)
	}
	if got := deriveUserName("Jane Doe", "ignored@example.com"); got != "Jane Doe" {
		t.Fatalf("expected explicit name, got %q", got)
	}
	if got := deriveUserName("", "invalid-email"); got != "Customer" {
		t.Fatalf("expected generic fallback, got %q", got)
	}
}
