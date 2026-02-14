package cryptox

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidPaymentLinkID    = errors.New("invalid payment link id")
	ErrInvalidPaymentLinkToken = errors.New("invalid payment link token")
	ErrExpiredPaymentLinkToken = errors.New("expired payment link token")
)

var paymentLinkIDPattern = regexp.MustCompile(`^pl_[a-f0-9]{32}$`)

type paymentLinkTokenClaims struct {
	LinkID    string `json:"lid"`
	ExpiresAt int64  `json:"exp"`
}

func IsValidPaymentLinkID(linkID string) bool {
	return paymentLinkIDPattern.MatchString(strings.TrimSpace(linkID))
}

func CreatePaymentLinkToken(c ICrypto, linkID string, expiresAt time.Time) (string, error) {
	if c == nil || !IsValidPaymentLinkID(linkID) {
		return "", ErrInvalidPaymentLinkID
	}
	claims := paymentLinkTokenClaims{LinkID: linkID, ExpiresAt: expiresAt.UTC().Unix()}
	if claims.ExpiresAt <= 0 {
		return "", ErrInvalidPaymentLinkToken
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	return c.Encrypt(string(payload))
}

func VerifyPaymentLinkToken(c ICrypto, token, expectedLinkID string, now time.Time) error {
	if c == nil || !IsValidPaymentLinkID(expectedLinkID) {
		return ErrInvalidPaymentLinkID
	}
	if strings.TrimSpace(token) == "" {
		return ErrInvalidPaymentLinkToken
	}
	decrypted, err := c.Decrypt(token)
	if err != nil {
		return ErrInvalidPaymentLinkToken
	}
	var claims paymentLinkTokenClaims
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		return ErrInvalidPaymentLinkToken
	}
	if !IsValidPaymentLinkID(claims.LinkID) || claims.LinkID != expectedLinkID {
		return ErrInvalidPaymentLinkToken
	}
	if claims.ExpiresAt <= 0 {
		return ErrInvalidPaymentLinkToken
	}
	if now.UTC().Unix() > claims.ExpiresAt {
		return ErrExpiredPaymentLinkToken
	}
	return nil
}
