package smtprelay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type postmarkEmail struct {
	From          string `json:"From"`
	To            string `json:"To"`
	Subject       string `json:"Subject"`
	HtmlBody      string `json:"HtmlBody,omitempty"`
	TextBody      string `json:"TextBody,omitempty"`
	MessageStream string `json:"MessageStream,omitempty"`
}

func (s *Session) sendViaPostmark() error {
	subject, htmlBody, textBody := parseEmailData(string(s.data))
	email := postmarkEmail{
		From: s.from, To: strings.Join(s.to, ","), Subject: subject,
		HtmlBody: htmlBody, TextBody: textBody, MessageStream: "outbound",
	}
	body, err := json.Marshal(email)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.backend.cfg.PostmarkAPI, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.backend.cfg.PostmarkToken)
	resp, err := s.backend.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("sent via Postmark status=%d to=%s", resp.StatusCode, strings.Join(s.to, ","))
		return nil
	}
	return fmt.Errorf("postmark API returned status %d: %s", resp.StatusCode, string(respBody))
}
