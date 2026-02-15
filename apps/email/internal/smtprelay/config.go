package smtprelay

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	SMTPAddr      string
	EmailProvider string
	PostmarkAPI   string
	PostmarkToken string
	ResendAPI     string
	ResendToken   string
}

func LoadConfigFromEnv() (Config, error) {
	cfg := Config{
		SMTPAddr:      envOrDefault("SMTP_ADDR", ":2525"),
		EmailProvider: strings.ToLower(envOrDefault("EMAIL_PROVIDER", "local")),
		PostmarkAPI:   envOrDefault("POSTMARK_API_URL", "https://api.postmarkapp.com/email"),
		PostmarkToken: os.Getenv("POSTMARK_API_KEY"),
		ResendAPI:     envOrDefault("RESEND_API_URL", "https://api.resend.com/emails"),
		ResendToken:   os.Getenv("RESEND_API_KEY"),
	}
	switch cfg.EmailProvider {
	case "local":
		return cfg, nil
	case "postmark":
		if cfg.PostmarkToken == "" {
			return Config{}, fmt.Errorf("POSTMARK_API_KEY is required when EMAIL_PROVIDER=postmark")
		}
	case "resend":
		if cfg.ResendToken == "" {
			return Config{}, fmt.Errorf("RESEND_API_KEY is required when EMAIL_PROVIDER=resend")
		}
	default:
		return Config{}, fmt.Errorf("unsupported EMAIL_PROVIDER=%q", cfg.EmailProvider)
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
