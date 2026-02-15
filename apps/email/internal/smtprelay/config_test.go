package smtprelay

import "testing"

func TestLoadConfigResendWithoutPostmarkKey(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "resend")
	t.Setenv("RESEND_API_KEY", "re_test")
	t.Setenv("POSTMARK_API_KEY", "")
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EmailProvider != "resend" {
		t.Fatalf("expected resend provider, got %q", cfg.EmailProvider)
	}
}

func TestLoadConfigPostmarkRequiresKey(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "postmark")
	t.Setenv("POSTMARK_API_KEY", "")
	if _, err := LoadConfigFromEnv(); err == nil {
		t.Fatal("expected error when POSTMARK_API_KEY is missing")
	}
}
