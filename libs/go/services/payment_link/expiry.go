package payment_link

import "time"

const (
	defaultLinkTTL = 7 * 24 * time.Hour
	maxLinkTTL     = 90 * 24 * time.Hour
)

func resolveLinkExpiry(expiresAt *time.Time, now time.Time) (time.Time, error) {
	if expiresAt == nil {
		return now.Add(defaultLinkTTL), nil
	}
	if !expiresAt.After(now) {
		return time.Time{}, NewValidationError("expires_at must be in the future")
	}
	if expiresAt.After(now.Add(maxLinkTTL)) {
		return time.Time{}, NewValidationError("expires_at cannot exceed 90 days")
	}
	return expiresAt.UTC(), nil
}
