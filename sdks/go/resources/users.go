package resources

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Address represents a billing address
type Address struct {
	ID         string    `json:"id"`
	Line1      string    `json:"line1"`
	Line2      string    `json:"line2,omitempty"`
	City       string    `json:"city"`
	State      string    `json:"state,omitempty"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// User represents a user in the system
type User struct {
	ID             string                 `json:"id"`
	Email          string                 `json:"email"`
	Name           *string                `json:"name,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	BillingAddress *Address               `json:"billing_address,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// ListUsersParams are the parameters for listing users
type ListUsersParams struct {
	Email  string
	Limit  int
	Offset int
}

// CreateUserRequest is the request to create a user
type CreateUserRequest struct {
	ID    *string `json:"id,omitempty"`
	Email string  `json:"email"`
	Name  *string `json:"name,omitempty"`
}

// Users provides access to user-related API operations
type Users struct {
	http *HTTPClient
}

// NewUsers creates a new Users resource
func NewUsers(http *HTTPClient) *Users {
	return &Users{http: http}
}

// List returns all users matching the given parameters
func (u *Users) List(ctx context.Context, params *ListUsersParams) ([]User, error) {
	query := url.Values{}
	if params != nil {
		if params.Email != "" {
			query.Set("email", params.Email)
		}
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.Offset > 0 {
			query.Set("offset", fmt.Sprintf("%d", params.Offset))
		}
	}

	path := "/v1/users"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var resp DataWrapper[[]User]
	if err := u.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a user by ID
func (u *Users) Get(ctx context.Context, userID string) (*User, error) {
	var resp DataWrapper[User]
	if err := u.http.Get(ctx, fmt.Sprintf("/v1/users/%s", userID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetByEmail returns a user by email, or nil if not found
func (u *Users) GetByEmail(ctx context.Context, email string) (*User, error) {
	users, err := u.List(ctx, &ListUsersParams{Email: email, Limit: 1})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

// Create creates a new user
func (u *Users) Create(ctx context.Context, data CreateUserRequest) (*User, error) {
	var resp DataWrapper[User]
	if err := u.http.Post(ctx, "/v1/users", data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// FindOrCreate finds a user by email or creates a new one if not found
func (u *Users) FindOrCreate(ctx context.Context, email string, name *string) (*User, error) {
	existing, err := u.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	// Create with default name if not provided
	if name == nil {
		defaultName := strings.Split(email, "@")[0]
		name = &defaultName
	}

	return u.Create(ctx, CreateUserRequest{
		Email: email,
		Name:  name,
	})
}

// UserSubscription represents a user's subscription with plan details
type UserSubscription struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	PlanID           *string    `json:"plan_id,omitempty"`
	Plan             *Plan      `json:"plan,omitempty"`
	ScheduledPlanID  *string    `json:"scheduled_plan_id,omitempty"`
	ScheduledPlan    *Plan      `json:"scheduled_plan,omitempty"`
	Status           string     `json:"status"`
	IsFree           bool       `json:"is_free"`
	HasPaymentMethod bool       `json:"has_payment_method"`
	NextResetAt      *time.Time `json:"next_reset_at,omitempty"`
}

// GetSubscription returns the active subscription for a user
func (u *Users) GetSubscription(ctx context.Context, userID string) (*UserSubscription, error) {
	var resp DataWrapper[UserSubscription]
	if err := u.http.Get(ctx, fmt.Sprintf("/v1/users/%s/subscription", userID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpsertBillingAddressInput is the input for upserting a billing address
type UpsertBillingAddressInput struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// UpsertBillingAddress creates or updates the billing address for a user
func (u *Users) UpsertBillingAddress(ctx context.Context, userID string, input UpsertBillingAddressInput) (*Address, error) {
	var resp DataWrapper[Address]
	if err := u.http.Post(ctx, fmt.Sprintf("/v1/users/%s/billing-address", userID), input, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
