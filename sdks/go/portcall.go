// Package portcall provides a Go SDK for the Portcall API.
//
// This SDK allows you to interact with the Portcall API for metered billing,
// entitlements, and feature management.
//
// Basic usage:
//
//	client := portcall.New(portcall.Config{
//		APIKey:  "sk_your_api_key",
//		BaseURL: "https://api.useportcall.com", // optional
//	})
//
//	// List plans
//	plans, err := client.Plans.List(ctx)
//
//	// Check entitlement
//	entitlement, err := client.Entitlements.Get(ctx, "user_123", "feature_id")
package portcall

import (
	"github.com/useportcall/portcall/sdks/go/resources"
)

const (
	// DefaultBaseURL is the default Portcall API endpoint
	DefaultBaseURL = "https://api.useportcall.com"
)

// Re-export types from resources package for convenience
type (
	Plan                       = resources.Plan
	PlanItem                   = resources.PlanItem
	PlanFeature                = resources.PlanFeature
	Subscription               = resources.Subscription
	User                       = resources.User
	Feature                    = resources.Feature
	Entitlement                = resources.Entitlement
	QuotaStatus                = resources.QuotaStatus
	MeterEvent                 = resources.MeterEvent
	CheckoutSession            = resources.CheckoutSession
	Invoice                    = resources.Invoice
	Address                    = resources.Address
	CreateSubscriptionInput    = resources.CreateSubscriptionRequest
	UpdateSubscriptionInput    = resources.UpdateSubscriptionRequest
	CreateUserInput            = resources.CreateUserRequest
	CreateMeterEventInput      = resources.CreateMeterEventRequest
	CreateCheckoutSessionInput = resources.CreateCheckoutSessionRequest
	CreateFeatureInput         = resources.CreateFeatureRequest
	CreatePlanInput            = resources.CreatePlanRequest
	UpdatePlanInput            = resources.UpdatePlanRequest
	AddFeatureInput            = resources.CreatePlanFeatureRequest
	CancelSubscriptionInput    = resources.CancelSubscriptionRequest
	UpsertBillingAddressInput  = resources.UpsertBillingAddressInput
	ListSubscriptionsParams    = resources.ListSubscriptionsParams
	ListUsersParams            = resources.ListUsersParams
	ListFeaturesParams         = resources.ListFeaturesParams
	ListInvoicesParams         = resources.ListInvoicesParams
	APIError                   = resources.APIError
)

// Config contains the configuration for the Portcall client
type Config struct {
	// APIKey is the API secret key (required)
	APIKey string

	// BaseURL is the base URL for the API (optional, defaults to production)
	BaseURL string
}

// Client is the main Portcall SDK client
type Client struct {
	config Config
	http   *resources.HTTPClient

	// Resource modules
	Plans            *resources.Plans
	Features         *resources.Features
	Subscriptions    *resources.Subscriptions
	Users            *resources.Users
	Entitlements     *resources.Entitlements
	MeterEvents      *resources.MeterEvents
	CheckoutSessions *resources.CheckoutSessions
	Invoices         *resources.Invoices
}

// New creates a new Portcall client
func New(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}

	http := resources.NewHTTPClient(config.APIKey, config.BaseURL)

	return &Client{
		config:           config,
		http:             http,
		Plans:            resources.NewPlans(http),
		Features:         resources.NewFeatures(http),
		Subscriptions:    resources.NewSubscriptions(http),
		Users:            resources.NewUsers(http),
		Entitlements:     resources.NewEntitlements(http),
		MeterEvents:      resources.NewMeterEvents(http),
		CheckoutSessions: resources.NewCheckoutSessions(http),
		Invoices:         resources.NewInvoices(http),
	}
}

// SetBaseURL changes the base URL for the client (useful for testing)
func (c *Client) SetBaseURL(url string) {
	c.config.BaseURL = url
	c.http.SetBaseURL(url)
}
