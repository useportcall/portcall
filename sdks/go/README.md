# Portcall Go SDK

A Go SDK for the Portcall API, providing metered billing, entitlements, and feature management.

## Installation

```bash
go get github.com/useportcall/portcall/sdks/go
```

## Code Generation CLI

The Go SDK includes a CLI tool that generates typed Go code from your Portcall app configuration, providing compile-time type safety for your plan and feature IDs.

### Generate Types

```bash
# Interactive mode (prompts for environment and API key)
go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen@latest generate

# With API key directly
go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen@latest generate --key sk_xxx

# Specify local environment
go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen@latest generate --url http://localhost:9080

# Custom output directory and package name
go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen@latest generate --output ./internal/billing --package billing
```

### Generated Code

This generates a `portcall/` directory (or your custom output) with:

- `generated.go` - Type-safe plan and feature ID constants
- `portcall.yaml` - YAML reference of your app configuration

### Using Generated Types

```go
package main

import (
    "context"
    "os"

    portcall "github.com/useportcall/portcall/sdks/go"
    "yourproject/portcall" // Import generated types
)

func main() {
    client := portcall.New(portcall.Config{
        APIKey: os.Getenv("PC_API_SECRET"),
    })

    ctx := context.Background()

    // Use type-safe feature constants instead of strings
    ent, err := client.Entitlements.Get(ctx, userID, billing.Features.MaxSubscriptions)

    // Use type-safe plan constants
    _, err = client.Subscriptions.Create(ctx, portcall.CreateSubscriptionInput{
        UserID: userID,
        PlanID: billing.Plans.Pro,
    })

    // Access plan metadata
    freePlan := billing.GetFreePlan()
    if freePlan != nil {
        fmt.Println("Free plan:", freePlan.Name)
    }
}
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    portcall "github.com/useportcall/portcall/sdks/go"
)

func main() {
    // Create a client
    client := portcall.New(portcall.Config{
        APIKey:  "sk_your_api_key",
        BaseURL: "https://api.useportcall.com", // optional, defaults to production
    })

    ctx := context.Background()

    // List plans
    plans, err := client.Plans.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d plans\n", len(plans))

    // Check entitlement
    entitlement, err := client.Entitlements.Get(ctx, "user_123", "max_subscriptions")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Usage: %d/%d\n", entitlement.Usage, entitlement.Quota)

    // Record meter event
    err = client.MeterEvents.Record(ctx, "user_123", "api_calls", 1)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Resources

### Plans

```go
// List all plans
plans, err := client.Plans.List(ctx)

// Get a specific plan
plan, err := client.Plans.Get(ctx, "plan_123")

// Create a plan
plan, err := client.Plans.Create(ctx, resources.CreatePlanRequest{
    Name:     "Pro Plan",
    Currency: "USD",
    Interval: "month",
})

// Publish a plan
plan, err := client.Plans.Publish(ctx, "plan_123")

// Get the free plan
freePlan, err := client.Plans.GetFreePlan(ctx)

// Get a pro plan
proPlan, err := client.Plans.GetProPlan(ctx)
```

### Subscriptions

```go
// List subscriptions for a user
subs, err := client.Subscriptions.List(ctx, &resources.ListSubscriptionsParams{
    UserID: "user_123",
    Status: "active",
})

// Create a subscription
sub, err := client.Subscriptions.Create(ctx, resources.CreateSubscriptionRequest{
    UserID: "user_123",
    PlanID: "plan_456",
})

// Cancel a subscription
sub, err := client.Subscriptions.Cancel(ctx, "sub_789", nil)

// Get active subscription for a user
activeSub, err := client.Subscriptions.GetActive(ctx, "user_123")
```

### Users

```go
// List users
users, err := client.Users.List(ctx, nil)

// Get a user by email
user, err := client.Users.GetByEmail(ctx, "user@example.com")

// Find or create a user
user, err := client.Users.FindOrCreate(ctx, "user@example.com", nil)
```

### Entitlements

```go
// Get an entitlement
ent, err := client.Entitlements.Get(ctx, "user_123", "max_subscriptions")

// Check if a feature is enabled
enabled, err := client.Entitlements.IsEnabled(ctx, "user_123", "premium_feature")

// Get remaining quota
remaining, err := client.Entitlements.GetRemaining(ctx, "user_123", "api_calls")

// Get full usage status
status, err := client.Entitlements.GetUsage(ctx, "user_123", "api_calls")
// status.Usage, status.Quota, status.Remaining, status.IsExceeded, status.IsUnlimited
```

### Meter Events

```go
// Record usage
err := client.MeterEvents.Record(ctx, "user_123", "api_calls", 10)

// Increment by 1
err := client.MeterEvents.Increment(ctx, "user_123", "api_calls")
```

### Checkout Sessions

```go
// Create a checkout session
session, err := client.CheckoutSessions.Create(ctx, resources.CreateCheckoutSessionRequest{
    UserID:      "user_123",
    PlanID:      "plan_456",
    CancelURL:   "https://example.com/cancel",
    RedirectURL: "https://example.com/success",
})
// Redirect user to session.URL
```

### Features

```go
// List all features
features, err := client.Features.List(ctx, nil)

// List only metered features
metered, err := client.Features.ListMetered(ctx)

// List only standard (boolean) features
standard, err := client.Features.ListStandard(ctx)
```

## Error Handling

The SDK returns errors wrapped in `*resources.APIError` for API errors:

```go
plan, err := client.Plans.Get(ctx, "invalid_id")
if err != nil {
    if apiErr, ok := err.(*resources.APIError); ok {
        fmt.Printf("API Error: %s (status: %d)\n", apiErr.Message, apiErr.Status)
    }
}
```

## Configuration

### Environment Variables

For production use, store your API key in an environment variable:

```go
client := portcall.New(portcall.Config{
    APIKey: os.Getenv("PORTCALL_API_SECRET"),
})
```

### Custom Base URL

For local development or testing:

```go
client := portcall.New(portcall.Config{
    APIKey:  "sk_test_...",
    BaseURL: "http://localhost:9080",
})
```

## License

Apache 2.0
