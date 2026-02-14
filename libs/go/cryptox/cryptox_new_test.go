package cryptox

import "testing"

func TestNewFromEnvMissingKey(t *testing.T) {
	t.Setenv("AES_ENCRYPTION_KEY", "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatalf("expected error when AES_ENCRYPTION_KEY is missing")
	}
}

func TestNewFromBase64KeyInvalid(t *testing.T) {
	_, err := NewFromBase64Key("not-base64")
	if err == nil {
		t.Fatalf("expected error for invalid base64 key")
	}
}
