package emailx

import (
	"fmt"
	"os"
)

type IEmailClient interface {
	Send(content, subject, from string, to []string) error
}

func New() (IEmailClient, error) { return NewFromEnv() }

func NewFromEnv() (IEmailClient, error) {
	emailProvider := os.Getenv("EMAIL_PROVIDER")
	switch emailProvider {
	case "postmark":
		apiKey := os.Getenv("POSTMARK_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("POSTMARK_API_KEY environment variable is not set")
		}

		return &postmarkEmailClient{apiKey: apiKey}, nil
	case "resend":
		apiKey := os.Getenv("RESEND_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("RESEND_API_KEY environment variable is not set")
		}
		return &resendEmailClient{apiKey: apiKey}, nil
	default:
		addr := os.Getenv("SMTP_SERVER")
		if addr == "" {
			return nil, fmt.Errorf("SMTP_SERVER environment variable is not set")
		}

		return &localEmailClient{addr: addr}, nil
	}
}
