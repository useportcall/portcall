package webhookx

import (
	"bytes"
	"fmt"
	"github.com/useportcall/portcall/libs/go/routerx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStripeWebhook_RelevantEventsAreEnqueued(t *testing.T) {
	for _, eventType := range []string{"payment_intent.payment_failed", "charge.failed", "invoice.payment_failed"} {
		t.Run(eventType, func(t *testing.T) {
			secret := "whsec_test"
			payload := []byte(fmt.Sprintf(`{"id":"evt_1","object":"event","type":"%s","data":{"object":{"id":"obj_1"}}}`, eventType))
			req := httptest.NewRequest(http.MethodPost, "/stripe/conn_1", bytes.NewReader(payload))
			req.Header.Set("Stripe-Signature", stripeSignatureHeader(secret, payload))

			rec := &queueRecorder{}
			r := routerx.New(&dbStub{conn: connectionWithSecret("conn_1", "enc_secret")}, &cryptoStub{decrypted: secret}, rec)
			RegisterRoutes(r)
			w := httptest.NewRecorder()
			r.Handler().ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d", w.Code)
			}
			if rec.count != 1 {
				t.Fatalf("expected enqueue=1, got %d", rec.count)
			}
		})
	}
}

func TestStripeWebhook_InvalidOrIrrelevantIgnored(t *testing.T) {
	secret := "whsec_test"
	payload := []byte(`{"id":"evt_2","object":"event","type":"customer.created","data":{"object":{"id":"obj_2"}}}`)
	req := httptest.NewRequest(http.MethodPost, "/stripe/conn_1", bytes.NewReader(payload))
	req.Header.Set("Stripe-Signature", stripeSignatureHeader(secret, payload))

	rec := &queueRecorder{}
	r := routerx.New(&dbStub{conn: connectionWithSecret("conn_1", "enc_secret")}, &cryptoStub{decrypted: secret}, rec)
	RegisterRoutes(r)
	w := httptest.NewRecorder()
	r.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if rec.count != 0 {
		t.Fatalf("expected enqueue=0, got %d", rec.count)
	}
}
