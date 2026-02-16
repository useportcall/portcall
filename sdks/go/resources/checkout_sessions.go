package resources

import (
	"context"
	"time"
)

// CheckoutSession represents a checkout session
type CheckoutSession struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	PlanID       string    `json:"plan_id"`
	Status       string    `json:"status"`
	URL          string    `json:"url,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateCheckoutSessionRequest is the request to create a checkout session
type CreateCheckoutSessionRequest struct {
	PlanID      string `json:"plan_id"`
	UserID      string `json:"user_id"`
	CancelURL   string `json:"cancel_url"`
	RedirectURL string `json:"redirect_url"`
}

// CheckoutSessions provides access to checkout session-related API operations
type CheckoutSessions struct {
	http *HTTPClient
}

// NewCheckoutSessions creates a new CheckoutSessions resource
func NewCheckoutSessions(http *HTTPClient) *CheckoutSessions {
	return &CheckoutSessions{http: http}
}

// Create creates a new checkout session
func (c *CheckoutSessions) Create(ctx context.Context, data CreateCheckoutSessionRequest) (*CheckoutSession, error) {
	var resp DataWrapper[CheckoutSession]
	if err := c.http.Post(ctx, "/v1/checkout-sessions", data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
