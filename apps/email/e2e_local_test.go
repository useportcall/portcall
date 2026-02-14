//go:build e2e_email_local

package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestEmailWorkerLocalE2E(t *testing.T) {
	redis, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer redis.Close()

	subjects := make(chan string, 2)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	api := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		_ = json.NewDecoder(r.Body).Decode(&payload)
		if s, ok := payload["subject"].(string); ok {
			subjects <- s
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"ok"}`))
	})}
	go func() { _ = api.Serve(ln) }()
	defer func() { _ = api.Close(); _ = ln.Close() }()

	t.Setenv("REDIS_ADDR", redis.Addr())
	t.Setenv("EMAIL_PROVIDER", "resend")
	t.Setenv("RESEND_API_KEY", "re_test")
	t.Setenv("RESEND_API_URL", "http://"+ln.Addr().String())
	t.Setenv("EMAIL_FROM", "relay-test@mail.useportcall.com")
	if _, err := startEmailWorker(); err != nil {
		t.Fatal(err)
	}

	if err := enqueueInvoicePaid("hello@useportcall.com"); err != nil {
		t.Fatal(err)
	}
	if err := enqueueStatus("hello@useportcall.com", "Local E2E Invoice Issued"); err != nil {
		t.Fatal(err)
	}

	deadline := time.After(8 * time.Second)
	got := []string{}
	for len(got) < 2 {
		select {
		case s := <-subjects:
			got = append(got, s)
		case <-deadline:
			t.Fatalf("timed out waiting for emails, got=%v", got)
		}
	}
	if !(contains(got, "Invoice Paid") && contains(got, "Local E2E Invoice Issued")) {
		t.Fatalf("unexpected subjects: %v", got)
	}
}

func contains(v []string, want string) bool {
	for _, s := range v {
		if strings.Contains(s, want) {
			return true
		}
	}
	return false
}
