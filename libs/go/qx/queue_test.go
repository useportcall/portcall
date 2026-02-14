package qx

import "testing"

func TestNewFromEnvMissingRedisAddr(t *testing.T) {
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_TLS", "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatalf("expected error when REDIS_ADDR is missing")
	}
}
