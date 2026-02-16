package resources

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
)

// Entitlement represents a user's entitlement to a feature
type Entitlement struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id,omitempty"`
	FeatureID   string     `json:"feature_id,omitempty"`
	Enabled     bool       `json:"enabled"`
	Value       any        `json:"value,omitempty"`
	Limit       *int64     `json:"limit,omitempty"`
	Quota       int64      `json:"quota"`
	Usage       int64      `json:"usage"`
	Interval    string     `json:"interval,omitempty"`
	ResetPeriod *string    `json:"reset_period,omitempty"`
	NextResetAt *time.Time `json:"next_reset_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// QuotaStatus represents the usage and quota status for a feature
type QuotaStatus struct {
	FeatureID   string `json:"feature_id"`
	Usage       int64  `json:"usage"`
	Quota       int64  `json:"quota"`
	Remaining   int64  `json:"remaining"`
	IsExceeded  bool   `json:"is_exceeded"`
	IsUnlimited bool   `json:"is_unlimited"`
}

// Entitlements provides access to entitlement-related API operations
type Entitlements struct {
	http  *HTTPClient
	users *Users
}

// NewEntitlements creates a new Entitlements resource
func NewEntitlements(http *HTTPClient) *Entitlements {
	return &Entitlements{
		http:  http,
		users: NewUsers(http),
	}
}

// Get returns a user's entitlement for a specific feature
func (e *Entitlements) Get(ctx context.Context, userID, featureID string) (*Entitlement, error) {
	// Resolve email to user ID if needed
	resolvedUserID, err := e.resolveUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var resp DataWrapper[Entitlement]
	if err := e.http.Get(ctx, fmt.Sprintf("/v1/entitlements/%s/%s", resolvedUserID, featureID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// IsEnabled checks if a user has an entitlement enabled
func (e *Entitlements) IsEnabled(ctx context.Context, userID, featureID string) (bool, error) {
	ent, err := e.Get(ctx, userID, featureID)
	if err != nil {
		return false, nil // Return false on error, not an error
	}
	return ent.Enabled, nil
}

// CanUse is an alias for IsEnabled
func (e *Entitlements) CanUse(ctx context.Context, userID, featureID string) (bool, error) {
	return e.IsEnabled(ctx, userID, featureID)
}

// HasAccess is an alias for IsEnabled
func (e *Entitlements) HasAccess(ctx context.Context, userID, featureID string) (bool, error) {
	return e.IsEnabled(ctx, userID, featureID)
}

// GetRemaining returns the remaining quota for a metered entitlement
// Returns -1 for unlimited, 0 on error or if not found
func (e *Entitlements) GetRemaining(ctx context.Context, userID, featureID string) (int64, error) {
	ent, err := e.Get(ctx, userID, featureID)
	if err != nil {
		return 0, nil // Return 0 on error
	}

	if ent.Quota < 0 {
		return -1, nil // Unlimited
	}

	remaining := ent.Quota - ent.Usage
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// GetUsage returns usage information for a metered feature
func (e *Entitlements) GetUsage(ctx context.Context, userID, featureID string) (*QuotaStatus, error) {
	ent, err := e.Get(ctx, userID, featureID)
	if err != nil {
		return nil, err
	}

	var remaining int64
	var isUnlimited bool
	if ent.Quota < 0 {
		remaining = math.MaxInt64
		isUnlimited = true
	} else {
		remaining = ent.Quota - ent.Usage
		if remaining < 0 {
			remaining = 0
		}
	}

	return &QuotaStatus{
		FeatureID:   featureID,
		Usage:       ent.Usage,
		Quota:       ent.Quota,
		Remaining:   remaining,
		IsExceeded:  !isUnlimited && ent.Usage >= ent.Quota,
		IsUnlimited: isUnlimited,
	}, nil
}

// CheckQuota checks if the user has remaining quota for a feature
func (e *Entitlements) CheckQuota(ctx context.Context, userID, featureID string) (*QuotaStatus, error) {
	return e.GetUsage(ctx, userID, featureID)
}

// resolveUserID resolves an email to a user ID if needed
func (e *Entitlements) resolveUserID(ctx context.Context, userIDOrEmail string) (string, error) {
	if strings.Contains(userIDOrEmail, "@") {
		user, err := e.users.GetByEmail(ctx, userIDOrEmail)
		if err != nil {
			return "", err
		}
		if user == nil {
			return "", &APIError{Message: fmt.Sprintf("user not found: %s", userIDOrEmail)}
		}
		return user.ID, nil
	}
	return userIDOrEmail, nil
}
