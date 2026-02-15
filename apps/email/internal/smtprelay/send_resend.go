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

func (s *Session) sendViaResend() error {
	subject, htmlBody, textBody := parseEmailData(string(s.data))
	payload := map[string]any{"from": s.from, "to": s.to, "subject": subject}
	if htmlBody != "" {
		payload["html"] = htmlBody
	}
	if textBody != "" {
		payload["text"] = textBody
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.backend.cfg.ResendAPI, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.backend.cfg.ResendToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.backend.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("sent via Resend status=%d to=%s", resp.StatusCode, strings.Join(s.to, ","))
		return nil
	}
	return fmt.Errorf("resend API returned status %d: %s", resp.StatusCode, string(respBody))
}
