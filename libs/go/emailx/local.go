package emailx

import (
	"fmt"
	"net/smtp"
)

type localEmailClient struct {
	addr string
}

func (c *localEmailClient) Send(content, subject, from string, to []string) error {
	if len(to) == 0 {
		return fmt.Errorf("recipient list is empty")
	}

	msg := []byte("" +
		"From: " + from + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		content +
		"\r\n")

	// No auth, no TLS; MailHog/MailDev accept this by default
	if err := smtp.SendMail(c.addr, nil, from, to, msg); err != nil {
		return err
	}

	return nil
}
