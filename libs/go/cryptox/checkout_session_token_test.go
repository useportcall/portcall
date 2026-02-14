package cryptox

import (
	"encoding/base64"
	"testing"
	"time"
)

func newTestCrypto(t *testing.T) ICrypto {
	t.Helper()
	key := base64.StdEncoding.EncodeToString([]byte("01234567890123456789012345678901"))
	t.Setenv("AES_ENCRYPTION_KEY", key)
	crypto, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return crypto
}

func TestCheckoutSessionToken_RoundTrip(t *testing.T) {
	c := newTestCrypto(t)
	sessionID := "cs_1234567890abcdef1234567890abcdef"
	expiresAt := time.Now().Add(10 * time.Minute)

	token, err := CreateCheckoutSessionToken(c, sessionID, expiresAt)
	if err != nil {
		t.Fatalf("CreateCheckoutSessionToken() error = %v", err)
	}

	if err := VerifyCheckoutSessionToken(c, token, sessionID, time.Now()); err != nil {
		t.Fatalf("VerifyCheckoutSessionToken() error = %v", err)
	}
}

func TestCheckoutSessionToken_RejectsExpiredToken(t *testing.T) {
	c := newTestCrypto(t)
	sessionID := "cs_1234567890abcdef1234567890abcdef"
	expiresAt := time.Now().Add(-1 * time.Minute)

	token, err := CreateCheckoutSessionToken(c, sessionID, expiresAt)
	if err != nil {
		t.Fatalf("CreateCheckoutSessionToken() error = %v", err)
	}

	err = VerifyCheckoutSessionToken(c, token, sessionID, time.Now())
	if err != ErrExpiredCheckoutSessionToken {
		t.Fatalf("expected ErrExpiredCheckoutSessionToken, got %v", err)
	}
}

func TestCheckoutSessionToken_RejectsWrongSessionID(t *testing.T) {
	c := newTestCrypto(t)
	sessionID := "cs_1234567890abcdef1234567890abcdef"

	token, err := CreateCheckoutSessionToken(c, sessionID, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("CreateCheckoutSessionToken() error = %v", err)
	}

	err = VerifyCheckoutSessionToken(c, token, "cs_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", time.Now())
	if err != ErrInvalidCheckoutSessionToken {
		t.Fatalf("expected ErrInvalidCheckoutSessionToken, got %v", err)
	}
}

func TestCheckoutSessionToken_RejectsInvalidIDs(t *testing.T) {
	c := newTestCrypto(t)
	_, err := CreateCheckoutSessionToken(c, "bad-id", time.Now().Add(time.Hour))
	if err != ErrInvalidCheckoutSessionID {
		t.Fatalf("expected ErrInvalidCheckoutSessionID, got %v", err)
	}

	err = VerifyCheckoutSessionToken(c, "x", "bad-id", time.Now())
	if err != ErrInvalidCheckoutSessionID {
		t.Fatalf("expected ErrInvalidCheckoutSessionID, got %v", err)
	}
}
