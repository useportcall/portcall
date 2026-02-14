package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"html/template"

	"github.com/useportcall/portcall/libs/go/qx/server"
)

func SendQuoteEmail(c server.IContext) error {
	var payload map[string]any
	if err := json.Unmarshal(c.Payload(), &payload); err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("templates/quote_notification.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var htmlContentBuf bytes.Buffer
	if err := tmpl.Execute(&htmlContentBuf, payload); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	content := htmlContentBuf.String()

	// Get sender email from environment variable
	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		return fmt.Errorf("EMAIL_FROM environment variable is not set")
	}

	// Get recipient email from payload
	to, ok := payload["recipient_email"].(string)
	if !ok || to == "" {
		return fmt.Errorf("recipient_email not found in payload")
	}

	return c.EmailClient().Send(
		content,
		"Quote Issued",
		from,
		[]string{to},
	)
}
