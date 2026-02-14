package billing

//go:generate go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen generate --output . --package billing

import (
	"context"
	"errors"
	"os"
	"sync"

	portcall "github.com/useportcall/portcall/sdks/go"
)

var (
	ErrNotConfigured       = errors.New("billing: not configured")
	ErrEntitlementNotFound = errors.New("billing: entitlement not found")
)

var (
	clientOnce sync.Once
	client     *portcall.Client
)

func GetClient() (*portcall.Client, error) {
	var initErr error

	clientOnce.Do(func() {
		apiKey := os.Getenv("DOGFOOD_LIVE_SECRET")
		if apiKey == "" {
			apiKey = os.Getenv("DOGFOOD_TEST_SECRET")
		}
		if apiKey == "" {
			initErr = ErrNotConfigured
			return
		}

		baseURL := os.Getenv("PORTCALL_API_URL")
		if baseURL == "" {
			baseURL = "https://api.useportcall.com"
		}

		client = portcall.New(portcall.Config{
			APIKey:  apiKey,
			BaseURL: baseURL,
		})
	})

	if initErr != nil {
		return nil, initErr
	}
	if client == nil {
		return nil, ErrNotConfigured
	}

	return client, nil
}

func IsConfigured() bool {
	_, err := GetClient()
	return err == nil
}

func CheckEntitlement(ctx context.Context, userID, featureID string) (*portcall.Entitlement, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	ent, err := c.Entitlements.Get(ctx, userID, featureID)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

func RecordUsage(ctx context.Context, userID, featureID string, quantity int) error {
	c, err := GetClient()
	if err != nil {
		return err
	}

	return c.MeterEvents.Record(ctx, userID, featureID, int64(quantity))
}

func CreateUser(ctx context.Context, email, name string) (*portcall.User, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	return c.Users.FindOrCreate(ctx, email, &name)
}

func CreateSubscription(ctx context.Context, userID, planID string) (*portcall.Subscription, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	return c.Subscriptions.Create(ctx, portcall.CreateSubscriptionInput{
		UserID: userID,
		PlanID: planID,
	})
}

func GetFreePlan(ctx context.Context) (*portcall.Plan, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	return c.Plans.GetFreePlan(ctx)
}

func GetProPlan(ctx context.Context) (*portcall.Plan, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	return c.Plans.GetProPlan(ctx)
}

func CreateCheckoutSession(ctx context.Context, userID, planID, successURL, cancelURL string) (*portcall.CheckoutSession, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	return c.CheckoutSessions.Create(ctx, portcall.CreateCheckoutSessionInput{
		UserID:      userID,
		PlanID:      planID,
		RedirectURL: successURL,
		CancelURL:   cancelURL,
	})
}

func GetUser(ctx context.Context, userID string) (*portcall.User, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	return c.Users.Get(ctx, userID)
}

func UpsertBillingAddress(ctx context.Context, userID string, input portcall.UpsertBillingAddressInput) (*portcall.Address, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	return c.Users.UpsertBillingAddress(ctx, userID, input)
}
