//go:build e2e_live_email

package main

import (
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestEmailWorkerLiveResend(t *testing.T) {
	apiKey := os.Getenv("RESEND_API_KEY")
	recipient := os.Getenv("E2E_EMAIL_TO")
	from := os.Getenv("E2E_EMAIL_FROM")
	if apiKey == "" || recipient == "" || from == "" {
		t.Fatal("set RESEND_API_KEY, E2E_EMAIL_TO, and E2E_EMAIL_FROM")
	}

	redis, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer redis.Close()

	t.Setenv("REDIS_ADDR", redis.Addr())
	t.Setenv("EMAIL_PROVIDER", "resend")
	t.Setenv("RESEND_API_KEY", apiKey)
	t.Setenv("EMAIL_FROM", from)
	if _, err := startEmailWorker(); err != nil {
		t.Fatal(err)
	}

	subject := "Live E2E Invoice Issued " + time.Now().UTC().Format("20060102-150405")
	if err := enqueueStatus(recipient, subject); err != nil {
		t.Fatal(err)
	}
	if err := enqueueInvoicePaid(recipient); err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	t.Logf("check inbox for subjects: %s and Invoice Paid", subject)
}
