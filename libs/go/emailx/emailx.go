package emailx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

func SendTemplateEmail(payload []byte, tmplFileName, subject, from string, to []string) error {
	if len(to) == 0 {
		return fmt.Errorf("recipient list is empty")
	}
	var body any
	if err := json.Unmarshal(payload, &body); err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(tmplFileName)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var htmlContentBuf bytes.Buffer
	if err := tmpl.Execute(&htmlContentBuf, body); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	htmlContent := htmlContentBuf.String()

	msg := []byte("" +
		"From: " + from + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		htmlContent +
		"\r\n")

	// No auth, no TLS; MailHog/MailDev accept this by default
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		return fmt.Errorf("SMTP_SERVER environment variable is not set")
	}
	if err := smtp.SendMail(smtpServer, nil, from, to, msg); err != nil {
		return err
	}

	log.Println("Email sent")

	return nil
}
