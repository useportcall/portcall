package smtprelay

import (
	"io"
	"mime"
	"mime/multipart"
	"strings"
)

func parseMultipart(body io.Reader, contentType string) (htmlBody, textBody string) {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", ""
	}
	reader := multipart.NewReader(body, params["boundary"])
	for {
		part, err := reader.NextPart()
		if err != nil {
			break
		}
		data, err := io.ReadAll(part)
		if err != nil {
			continue
		}
		content := decodeBody(string(data), part.Header.Get("Content-Transfer-Encoding"))
		partType := part.Header.Get("Content-Type")
		if strings.Contains(partType, "text/html") {
			htmlBody = content
		}
		if strings.Contains(partType, "text/plain") {
			textBody = content
		}
	}
	return htmlBody, textBody
}

func decodeBody(body, encoding string) string {
	if strings.EqualFold(encoding, "quoted-printable") {
		body = strings.ReplaceAll(body, "=\r\n", "")
		body = strings.ReplaceAll(body, "=\n", "")
	}
	return body
}
