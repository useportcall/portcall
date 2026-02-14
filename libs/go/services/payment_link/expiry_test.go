package payment_link

import (
	"testing"
	"time"
)

func TestResolveLinkExpiry_Default(t *testing.T) {
	now := time.Now().UTC()
	expiresAt, err := resolveLinkExpiry(nil, now)
	if err != nil {
		t.Fatalf("resolveLinkExpiry() error = %v", err)
	}
	if expiresAt.Before(now.Add((7 * 24 * time.Hour) - time.Minute)) {
		t.Fatalf("expected default expiry around 7 days, got %s", expiresAt)
	}
}

func TestResolveLinkExpiry_RejectsPast(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-time.Minute)
	_, err := resolveLinkExpiry(&past, now)
	if err == nil {
		t.Fatal("expected validation error for past expiry")
	}
}
