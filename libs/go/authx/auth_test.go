package authx

import "testing"

func TestNewFromEnvKeycloakMissingURL(t *testing.T) {
	t.Setenv("AUTH_MODULE", "")
	t.Setenv("KEYCLOAK_API_URL", "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatalf("expected error when KEYCLOAK_API_URL is missing")
	}
}

func TestNewFromEnvClerkMissingSecret(t *testing.T) {
	t.Setenv("AUTH_MODULE", "clerk")
	t.Setenv("CLERK_SECRET", "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatalf("expected error when CLERK_SECRET is missing")
	}
}
