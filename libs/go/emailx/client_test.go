package emailx

import "testing"

func TestNewResendClient(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "resend")
	t.Setenv("RESEND_API_KEY", "re_test")
	client, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, ok := client.(*resendEmailClient); !ok {
		t.Fatalf("expected resendEmailClient, got %T", client)
	}
}

func TestNewLocalClient(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "local")
	t.Setenv("SMTP_SERVER", "localhost:1025")
	client, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, ok := client.(*localEmailClient); !ok {
		t.Fatalf("expected localEmailClient, got %T", client)
	}
}

func TestNewFromEnvMissingResendAPIKey(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "resend")
	t.Setenv("RESEND_API_KEY", "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatalf("expected error when RESEND_API_KEY is missing")
	}
}

func TestNewFromEnvMissingSMTPServer(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "")
	t.Setenv("SMTP_SERVER", "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatalf("expected error when SMTP_SERVER is missing")
	}
}
