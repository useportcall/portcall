package cryptox

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidCheckoutSessionID    = errors.New("invalid checkout session id")
	ErrInvalidCheckoutSessionToken = errors.New("invalid checkout session token")
	ErrExpiredCheckoutSessionToken = errors.New("expired checkout session token")
)

var checkoutSessionIDPattern = regexp.MustCompile(`^cs_[a-f0-9]{32}$`)

type checkoutSessionTokenClaims struct {
	SessionID string `json:"sid"`
	ExpiresAt int64  `json:"exp"`
}

func IsValidCheckoutSessionID(sessionID string) bool {
	return checkoutSessionIDPattern.MatchString(strings.TrimSpace(sessionID))
}

func CreateCheckoutSessionToken(c ICrypto, sessionID string, expiresAt time.Time) (string, error) {
	if c == nil || !IsValidCheckoutSessionID(sessionID) {
		return "", ErrInvalidCheckoutSessionID
	}

	claims := checkoutSessionTokenClaims{
		SessionID: sessionID,
		ExpiresAt: expiresAt.UTC().Unix(),
	}
	if claims.ExpiresAt <= 0 {
		return "", ErrInvalidCheckoutSessionToken
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	token, err := c.Encrypt(string(payload))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyCheckoutSessionToken(c ICrypto, token, expectedSessionID string, now time.Time) error {
	if c == nil || !IsValidCheckoutSessionID(expectedSessionID) {
		return ErrInvalidCheckoutSessionID
	}
	if strings.TrimSpace(token) == "" {
		return ErrInvalidCheckoutSessionToken
	}

	decrypted, err := c.Decrypt(token)
	if err != nil {
		return ErrInvalidCheckoutSessionToken
	}

	var claims checkoutSessionTokenClaims
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		return ErrInvalidCheckoutSessionToken
	}

	if !IsValidCheckoutSessionID(claims.SessionID) || claims.SessionID != expectedSessionID {
		return ErrInvalidCheckoutSessionToken
	}
	if claims.ExpiresAt <= 0 {
		return ErrInvalidCheckoutSessionToken
	}
	if now.UTC().Unix() > claims.ExpiresAt {
		return ErrExpiredCheckoutSessionToken
	}

	return nil
}
