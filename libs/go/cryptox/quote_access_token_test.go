package cryptox

import (
	"testing"
	"time"
)

func TestQuoteAccessToken_RoundTrip(t *testing.T) {
	c := newTestCrypto(t)
	quoteID := "quote_1234567890abcdef1234567890abcdef"
	expiresAt := time.Now().Add(10 * time.Minute)

	token, err := CreateQuoteAccessToken(c, quoteID, expiresAt)
	if err != nil {
		t.Fatalf("CreateQuoteAccessToken() error = %v", err)
	}

	if err := VerifyQuoteAccessToken(c, token, quoteID, time.Now()); err != nil {
		t.Fatalf("VerifyQuoteAccessToken() error = %v", err)
	}
}

func TestQuoteAccessToken_RejectsExpiredToken(t *testing.T) {
	c := newTestCrypto(t)
	quoteID := "quote_1234567890abcdef1234567890abcdef"
	expiresAt := time.Now().Add(-1 * time.Minute)

	token, err := CreateQuoteAccessToken(c, quoteID, expiresAt)
	if err != nil {
		t.Fatalf("CreateQuoteAccessToken() error = %v", err)
	}

	err = VerifyQuoteAccessToken(c, token, quoteID, time.Now())
	if err != ErrExpiredQuoteToken {
		t.Fatalf("expected ErrExpiredQuoteToken, got %v", err)
	}
}

func TestQuoteAccessToken_RejectsInvalidIDs(t *testing.T) {
	c := newTestCrypto(t)
	_, err := CreateQuoteAccessToken(c, "bad-id", time.Now().Add(time.Hour))
	if err != ErrInvalidQuoteID {
		t.Fatalf("expected ErrInvalidQuoteID, got %v", err)
	}

	err = VerifyQuoteAccessToken(c, "x", "bad-id", time.Now())
	if err != ErrInvalidQuoteID {
		t.Fatalf("expected ErrInvalidQuoteID, got %v", err)
	}
}
