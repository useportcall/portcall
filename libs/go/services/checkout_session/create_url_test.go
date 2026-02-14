package checkout_session

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type tokenCrypto struct {
	token string
}

func (m tokenCrypto) Encrypt(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", errors.New("empty payload")
	}
	return m.token, nil
}

func (m tokenCrypto) Decrypt(data string) (string, error) { return data, nil }
func (m tokenCrypto) CompareHash(hashed string, plain string) (bool, error) {
	return hashed == plain, nil
}

func TestBuildCheckoutURL_WithToken(t *testing.T) {
	session := &models.CheckoutSession{
		PublicID:  "cs_1234567890abcdef1234567890abcdef",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	got, err := buildCheckoutURL("https://checkout.example.com", session, tokenCrypto{token: "token123"})
	if err != nil {
		t.Fatalf("buildCheckoutURL() error = %v", err)
	}

	if !strings.Contains(got, "id=cs_1234567890abcdef1234567890abcdef") {
		t.Fatalf("expected id query param in URL, got %s", got)
	}
	if !strings.Contains(got, "st=token123") {
		t.Fatalf("expected st query param in URL, got %s", got)
	}
}

func TestBuildCheckoutURL_WithoutCrypto(t *testing.T) {
	session := &models.CheckoutSession{
		PublicID:  "cs_1234567890abcdef1234567890abcdef",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	got, err := buildCheckoutURL("https://checkout.example.com/path", session, nil)
	if err != nil {
		t.Fatalf("buildCheckoutURL() error = %v", err)
	}

	if !strings.Contains(got, "id=cs_1234567890abcdef1234567890abcdef") {
		t.Fatalf("expected id query param in URL, got %s", got)
	}
	if strings.Contains(got, "st=") {
		t.Fatalf("did not expect st query param in URL, got %s", got)
	}
}

func TestValidateCheckoutReturnURL(t *testing.T) {
	if _, err := validateCheckoutReturnURL("javascript:alert(1)", "redirect_url"); err == nil {
		t.Fatal("expected validation error for javascript URL")
	}

	got, err := validateCheckoutReturnURL("https://example.com/success", "redirect_url")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://example.com/success" {
		t.Fatalf("unexpected normalized URL %q", got)
	}
}
