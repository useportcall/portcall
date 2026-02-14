package payment_link

import (
	"net/url"
	"strings"
)

func validateReturnURL(raw string, field string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", NewValidationError("invalid %s", field)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return "", NewValidationError("invalid %s", field)
	}
	if parsed.Host == "" || parsed.User != nil {
		return "", NewValidationError("invalid %s", field)
	}
	return parsed.String(), nil
}
