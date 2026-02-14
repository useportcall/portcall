package payment_link

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func TestBuildPaymentLinkURL(t *testing.T) {
	key := base64.StdEncoding.EncodeToString([]byte("01234567890123456789012345678901"))
	t.Setenv("AES_ENCRYPTION_KEY", key)
	crypto, err := cryptox.New()
	if err != nil {
		t.Fatalf("cryptox.New() error = %v", err)
	}
	link := &models.PaymentLink{
		PublicID:  "pl_1234567890abcdef1234567890abcdef",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	url, err := buildPaymentLinkURL("https://checkout.example.com", link, crypto)
	if err != nil {
		t.Fatalf("buildPaymentLinkURL() error = %v", err)
	}
	if !strings.Contains(url, "pl=pl_1234567890abcdef1234567890abcdef") {
		t.Fatalf("expected payment link id in URL, got %s", url)
	}
	if !strings.Contains(url, "pt=") {
		t.Fatalf("expected payment link token in URL, got %s", url)
	}
}
