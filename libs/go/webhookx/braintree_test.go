package webhookx

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/braintree-go/braintree-go"
	"github.com/useportcall/portcall/libs/go/routerx"
)

const (
	btMerchantID = "mid"
	btPublicKey  = "sz9g7zhxz8838v7h"
	btPrivateKey = "0c809a2d2e8f4e4c817900ff441c9554"
)

func TestBraintreeChallenge_ReturnsVerificationString(t *testing.T) {
	secret := `{"merchant_id":"mid","private_key":"` + btPrivateKey + `","environment":"sandbox"}`
	rec := &queueRecorder{}
	r := routerx.New(&dbStub{conn: braintreeConnection("conn_1", "enc_key", btPublicKey)}, &cryptoStub{decrypted: secret}, rec)
	RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/braintree/conn_1?bt_challenge=abc123", nil)
	w := httptest.NewRecorder()
	r.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	gw := braintree.New(braintree.Sandbox, btMerchantID, btPublicKey, btPrivateKey)
	expected, _ := gw.WebhookNotification().Verify("abc123")
	if w.Body.String() != expected {
		t.Fatalf("unexpected challenge response: %s", w.Body.String())
	}
}

func TestBraintreeWebhook_RelevantEventEnqueued(t *testing.T) {
	secret := `{"merchant_id":"mid","private_key":"` + btPrivateKey + `","environment":"sandbox"}`
	gw := braintree.New(braintree.Sandbox, btMerchantID, btPublicKey, btPrivateKey)
	sampleReq, err := gw.WebhookTesting().Request(braintree.TransactionSettledWebhook, "tx_1")
	if err != nil {
		t.Fatalf("failed to create sample request: %v", err)
	}
	body, _ := io.ReadAll(sampleReq.Body)
	req := httptest.NewRequest(http.MethodPost, "/braintree/conn_1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := &queueRecorder{}
	r := routerx.New(&dbStub{conn: braintreeConnection("conn_1", "enc_key", btPublicKey)}, &cryptoStub{decrypted: secret}, rec)
	RegisterRoutes(r)
	w := httptest.NewRecorder()
	r.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if rec.count != 1 {
		t.Fatalf("expected enqueue=1, got %d", rec.count)
	}
}
