package emailx

import (
	"log"
	"os"
)

type IEmailClient interface {
	Send(content, subject, from string, to []string) error
}

func New() IEmailClient {
	emailProvider := os.Getenv("EMAIL_PROVIDER")
	switch emailProvider {
	case "postmark":
		apiKey := os.Getenv("POSTMARK_API_KEY")
		if apiKey == "" {
			log.Fatal("POSTMARK_API_KEY environment variable is not set")
		}

		return &postmarkEmailClient{apiKey: apiKey}
	default:
		addr := os.Getenv("SMTP_SERVER")
		if os.Getenv("SMTP_SERVER") == "" {
			log.Fatal("SMTP_SERVER environment variable is not set")
		}

		return &localEmailClient{addr: addr}
	}
}
