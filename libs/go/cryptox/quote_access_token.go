package cryptox

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidQuoteID    = errors.New("invalid quote id")
	ErrInvalidQuoteToken = errors.New("invalid quote token")
	ErrExpiredQuoteToken = errors.New("expired quote token")
)

var quoteIDPattern = regexp.MustCompile(`^quote_[a-f0-9]{32}$`)

type quoteAccessTokenClaims struct {
	QuoteID   string `json:"qid"`
	ExpiresAt int64  `json:"exp"`
}

func IsValidQuoteID(quoteID string) bool {
	return quoteIDPattern.MatchString(strings.TrimSpace(quoteID))
}

func CreateQuoteAccessToken(c ICrypto, quoteID string, expiresAt time.Time) (string, error) {
	if c == nil || !IsValidQuoteID(quoteID) {
		return "", ErrInvalidQuoteID
	}
	exp := expiresAt.UTC().Unix()
	if exp <= 0 {
		return "", ErrInvalidQuoteToken
	}

	payload, err := json.Marshal(quoteAccessTokenClaims{
		QuoteID:   quoteID,
		ExpiresAt: exp,
	})
	if err != nil {
		return "", err
	}

	return c.Encrypt(string(payload))
}

func VerifyQuoteAccessToken(c ICrypto, token, expectedQuoteID string, now time.Time) error {
	if c == nil || !IsValidQuoteID(expectedQuoteID) {
		return ErrInvalidQuoteID
	}
	if strings.TrimSpace(token) == "" {
		return ErrInvalidQuoteToken
	}

	decrypted, err := c.Decrypt(token)
	if err != nil {
		return ErrInvalidQuoteToken
	}

	var claims quoteAccessTokenClaims
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		return ErrInvalidQuoteToken
	}

	if !IsValidQuoteID(claims.QuoteID) || claims.QuoteID != expectedQuoteID {
		return ErrInvalidQuoteToken
	}
	if claims.ExpiresAt <= 0 {
		return ErrInvalidQuoteToken
	}
	if now.UTC().Unix() > claims.ExpiresAt {
		return ErrExpiredQuoteToken
	}

	return nil
}
