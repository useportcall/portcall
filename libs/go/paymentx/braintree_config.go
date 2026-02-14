package paymentx

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type BraintreeCredentials struct {
	MerchantID      string
	PrivateKey      string
	Environment     string
	MerchantAccount string
}

type braintreeSecretJSON struct {
	MerchantID      string `json:"merchant_id"`
	PrivateKey      string `json:"private_key"`
	Environment     string `json:"environment"`
	MerchantAccount string `json:"merchant_account_id"`
}

func ParseBraintreeCredentials(raw string) (*BraintreeCredentials, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty braintree secret")
	}
	if cfg := parseBraintreeJSON(raw); cfg != nil {
		return normalizeBraintreeCredentials(cfg)
	}
	parts := splitBraintreeSecret(raw)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid braintree secret format")
	}
	cfg := &BraintreeCredentials{
		MerchantID: parts[0],
		PrivateKey: parts[1],
	}
	if len(parts) > 2 {
		cfg.Environment = parts[2]
	}
	if len(parts) > 3 {
		cfg.MerchantAccount = parts[3]
	}
	return normalizeBraintreeCredentials(cfg)
}

func parseBraintreeJSON(raw string) *BraintreeCredentials {
	var payload braintreeSecretJSON
	if json.Unmarshal([]byte(raw), &payload) != nil {
		return nil
	}
	if payload.MerchantID == "" || payload.PrivateKey == "" {
		return nil
	}
	return &BraintreeCredentials{
		MerchantID:      payload.MerchantID,
		PrivateKey:      payload.PrivateKey,
		Environment:     payload.Environment,
		MerchantAccount: payload.MerchantAccount,
	}
}

func splitBraintreeSecret(raw string) []string {
	separator := ":"
	if strings.Contains(raw, "|") && !strings.Contains(raw, ":") {
		separator = "|"
	}
	parts := strings.SplitN(raw, separator, 4)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func normalizeBraintreeCredentials(cfg *BraintreeCredentials) (*BraintreeCredentials, error) {
	cfg.MerchantID = strings.TrimSpace(cfg.MerchantID)
	cfg.PrivateKey = strings.TrimSpace(cfg.PrivateKey)
	cfg.Environment = strings.ToLower(strings.TrimSpace(cfg.Environment))
	cfg.MerchantAccount = strings.TrimSpace(cfg.MerchantAccount)
	if cfg.MerchantID == "" || cfg.PrivateKey == "" {
		return nil, fmt.Errorf("missing braintree merchant_id/private_key")
	}
	if cfg.Environment == "" {
		cfg.Environment = strings.ToLower(strings.TrimSpace(os.Getenv("BRAINTREE_ENVIRONMENT")))
	}
	if cfg.Environment == "" {
		cfg.Environment = "production"
	}
	switch cfg.Environment {
	case "production", "sandbox", "development":
		return cfg, nil
	default:
		return nil, fmt.Errorf("unsupported braintree environment: %s", cfg.Environment)
	}
}
