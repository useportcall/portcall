package dbx

import "testing"

func TestNewFromEnvMissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatalf("expected error when DATABASE_URL is missing")
	}
}
