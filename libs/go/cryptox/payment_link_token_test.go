package cryptox

import (
	"testing"
	"time"
)

func TestPaymentLinkToken_RoundTrip(t *testing.T) {
	c := newTestCrypto(t)
	linkID := "pl_1234567890abcdef1234567890abcdef"
	expiresAt := time.Now().Add(time.Hour)
	token, err := CreatePaymentLinkToken(c, linkID, expiresAt)
	if err != nil {
		t.Fatalf("CreatePaymentLinkToken() error = %v", err)
	}
	if err := VerifyPaymentLinkToken(c, token, linkID, time.Now()); err != nil {
		t.Fatalf("VerifyPaymentLinkToken() error = %v", err)
	}
}

func TestPaymentLinkToken_RejectsExpiredToken(t *testing.T) {
	c := newTestCrypto(t)
	linkID := "pl_1234567890abcdef1234567890abcdef"
	expiresAt := time.Now().Add(-time.Minute)
	token, err := CreatePaymentLinkToken(c, linkID, expiresAt)
	if err != nil {
		t.Fatalf("CreatePaymentLinkToken() error = %v", err)
	}
	err = VerifyPaymentLinkToken(c, token, linkID, time.Now())
	if err != ErrExpiredPaymentLinkToken {
		t.Fatalf("expected ErrExpiredPaymentLinkToken, got %v", err)
	}
}

func TestPaymentLinkToken_RejectsWrongLinkID(t *testing.T) {
	c := newTestCrypto(t)
	linkID := "pl_1234567890abcdef1234567890abcdef"
	token, err := CreatePaymentLinkToken(c, linkID, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("CreatePaymentLinkToken() error = %v", err)
	}
	err = VerifyPaymentLinkToken(c, token, "pl_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", time.Now())
	if err != ErrInvalidPaymentLinkToken {
		t.Fatalf("expected ErrInvalidPaymentLinkToken, got %v", err)
	}
}
