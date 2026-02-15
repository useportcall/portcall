package smtprelay

import (
	"io"
	"mime"
	"net/mail"
	"regexp"
	"strings"
)

func parseEmailData(raw string) (subject, htmlBody, textBody string) {
	msg, err := mail.ReadMessage(strings.NewReader(raw))
	if err != nil {
		return parseEmailDataBasic(raw)
	}
	dec := new(mime.WordDecoder)
	subject, err = dec.DecodeHeader(msg.Header.Get("Subject"))
	if err != nil {
		subject = msg.Header.Get("Subject")
	}
	contentType := msg.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {
		htmlBody, textBody = parseMultipart(msg.Body, contentType)
		return subject, htmlBody, textBody
	}
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return subject, "", ""
	}
	content := string(body)
	if strings.Contains(contentType, "text/html") {
		return subject, content, ""
	}
	return subject, "", content
}

func parseEmailDataBasic(raw string) (subject, htmlBody, textBody string) {
	lines := strings.Split(raw, "\n")
	body := []string{}
	headerDone := false
	for _, line := range lines {
		if !headerDone {
			if strings.HasPrefix(line, "Subject:") {
				subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
			}
			if strings.TrimSpace(line) == "" {
				headerDone = true
			}
			continue
		}
		body = append(body, line)
	}
	content := strings.Join(body, "\n")
	if strings.Contains(content, "<html") || strings.Contains(content, "<body") {
		return subject, content, ""
	}
	return subject, "", content
}

func stripHTMLTags(v string) string {
	return regexp.MustCompile(`<[^>]*>`).ReplaceAllString(v, "")
}
