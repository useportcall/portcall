package resources

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MeterEvent represents a metered usage event
type MeterEvent struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	FeatureID string                 `json:"feature_id"`
	Usage     int64                  `json:"usage"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// CreateMeterEventRequest is the request to create a meter event
type CreateMeterEventRequest struct {
	UserID    string `json:"user_id"`
	FeatureID string `json:"feature_id"`
	Usage     int64  `json:"usage"`
}

// MeterEvents provides access to meter event-related API operations
type MeterEvents struct {
	http  *HTTPClient
	users *Users
}

// NewMeterEvents creates a new MeterEvents resource
func NewMeterEvents(http *HTTPClient) *MeterEvents {
	return &MeterEvents{
		http:  http,
		users: NewUsers(http),
	}
}

// Create creates a meter event (records usage)
func (m *MeterEvents) Create(ctx context.Context, data CreateMeterEventRequest) (*MeterEvent, error) {
	var resp DataWrapper[MeterEvent]
	if err := m.http.Post(ctx, "/v1/meter-events", data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Record records usage for a metered feature (convenience method)
// userIDOrEmail can be a user ID or an email address
func (m *MeterEvents) Record(ctx context.Context, userIDOrEmail, featureID string, usage int64) error {
	// Resolve email to user ID if needed
	userID, err := m.resolveUserID(ctx, userIDOrEmail)
	if err != nil {
		return err
	}

	_, err = m.Create(ctx, CreateMeterEventRequest{
		UserID:    userID,
		FeatureID: featureID,
		Usage:     usage,
	})
	return err
}

// Increment is an alias for Record with usage=1
func (m *MeterEvents) Increment(ctx context.Context, userIDOrEmail, featureID string) error {
	return m.Record(ctx, userIDOrEmail, featureID, 1)
}

// resolveUserID resolves an email to a user ID if needed
func (m *MeterEvents) resolveUserID(ctx context.Context, userIDOrEmail string) (string, error) {
	if strings.Contains(userIDOrEmail, "@") {
		user, err := m.users.GetByEmail(ctx, userIDOrEmail)
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
