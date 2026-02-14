package payment_link

import "testing"

func TestValidateReturnURL_AllowsHTTPAndHTTPS(t *testing.T) {
	if _, err := validateReturnURL("https://example.com/success", "redirect_url"); err != nil {
		t.Fatalf("expected https URL to pass, got %v", err)
	}
	if _, err := validateReturnURL("http://localhost:3000/cancel", "cancel_url"); err != nil {
		t.Fatalf("expected http URL to pass, got %v", err)
	}
}

func TestValidateReturnURL_RejectsBadScheme(t *testing.T) {
	if _, err := validateReturnURL("javascript:alert(1)", "redirect_url"); err == nil {
		t.Fatal("expected invalid scheme to fail")
	}
}
