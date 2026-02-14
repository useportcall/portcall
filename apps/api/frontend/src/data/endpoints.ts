export const apiEndpoints = [
  {
    id: "checkout",
    title: "Checkout Sessions",
    endpoints: [
      {
        id: "create-checkout-session",
        title: "Create Checkout Session",
        method: "POST",
        path: "/v1/checkout-sessions",
        description:
          "Create a new checkout session for a user to subscribe to a plan. Returns a session URL that can be used to redirect the user to the checkout page.",
        authentication: true,
        parameters: [
          {
            name: "user_id",
            type: "string",
            required: true,
            description: "The ID of the user creating the checkout session",
            location: "body" as const,
          },
          {
            name: "plan_id",
            type: "string",
            required: true,
            description: "The ID of the plan to subscribe to",
            location: "body" as const,
          },
          {
            name: "redirect_url",
            type: "string",
            required: false,
            description: "URL to redirect to after successful checkout",
            location: "body" as const,
          },
          {
            name: "cancel_url",
            type: "string",
            required: false,
            description: "URL to redirect to if checkout is cancelled",
            location: "body" as const,
          },
        ],
        requestBody: {
          example: `{
  "user_id": "usr_1234567890",
  "plan_id": "plan_basic_monthly",
  "redirect_url": "https://yourapp.com/success",
  "cancel_url": "https://yourapp.com/cancel"
}`,
        },
        response: {
          example: `{
  "id": "cs_1234567890",
  "user_id": "usr_1234567890",
  "plan_id": "plan_basic_monthly",
  "url": "https://checkout.example.com?id=cs_1234567890abcdef1234567890abcdef&st=token_abc123",
  "redirect_url": "https://yourapp.com/success",
  "cancel_url": "https://yourapp.com/cancel",
  "expires_at": "2024-01-17T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z",
  "external_provider": "stripe",
  "external_session_id": "cs_test_abc123",
  "external_client_secret": "cs_test_abc123_secret_xyz",
  "external_public_key": "pk_test_xyz"
}`,
        },
      },
    ],
  },
  {
    id: "subscriptions",
    title: "Subscriptions",
    endpoints: [
      {
        id: "list-subscriptions",
        title: "List Subscriptions",
        method: "GET",
        path: "/v1/subscriptions",
        description:
          "Retrieve a list of subscriptions. You can filter by user_id to get subscriptions for a specific user.",
        authentication: true,
        parameters: [
          {
            name: "user_id",
            type: "string",
            required: false,
            description: "Filter subscriptions by user ID",
            location: "query" as const,
          },
          {
            name: "status",
            type: "string",
            required: false,
            description:
              "Filter by subscription status (active, cancelled, expired)",
            location: "query" as const,
          },
          {
            name: "limit",
            type: "number",
            required: false,
            description:
              "Maximum number of results to return (default: 20, max: 100)",
            location: "query" as const,
          },
          {
            name: "offset",
            type: "number",
            required: false,
            description: "Number of results to skip (for pagination)",
            location: "query" as const,
          },
        ],
        response: {
          example: `{
  "data": [
    {
      "id": "sub_1234567890",
      "user_id": "usr_1234567890",
      "plan_id": "plan_basic_monthly",
      "status": "active",
      "current_period_start": "2024-01-01T00:00:00Z",
      "current_period_end": "2024-02-01T00:00:00Z",
      "cancel_at_period_end": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}`,
        },
      },
      {
        id: "update-subscription",
        title: "Update Subscription",
        method: "POST",
        path: "/v1/subscriptions/{subscription_id}",
        description:
          "Update an existing subscription. You can change the plan, billing cycle, or other subscription details.",
        authentication: true,
        parameters: [
          {
            name: "subscription_id",
            type: "string",
            required: true,
            description: "The ID of the subscription to update",
            location: "path" as const,
          },
          {
            name: "plan_id",
            type: "string",
            required: false,
            description: "Change the subscription to a different plan",
            location: "body" as const,
          },
          {
            name: "proration_behavior",
            type: "string",
            required: false,
            description: "How to handle proration (create_prorations, none)",
            location: "body" as const,
          },
        ],
        requestBody: {
          example: `{
  "plan_id": "plan_premium_monthly",
  "proration_behavior": "create_prorations"
}`,
        },
        response: {
          example: `{
  "id": "sub_1234567890",
  "user_id": "usr_1234567890",
  "plan_id": "plan_premium_monthly",
  "status": "active",
  "current_period_start": "2024-01-01T00:00:00Z",
  "current_period_end": "2024-02-01T00:00:00Z",
  "updated_at": "2024-01-17T10:30:00Z"
}`,
        },
      },
      {
        id: "cancel-subscription",
        title: "Cancel Subscription",
        method: "POST",
        path: "/v1/subscriptions/{subscription_id}/cancel",
        description:
          "Cancel a subscription. By default, the subscription will remain active until the end of the current billing period.",
        authentication: true,
        parameters: [
          {
            name: "subscription_id",
            type: "string",
            required: true,
            description: "The ID of the subscription to cancel",
            location: "path" as const,
          },
          {
            name: "cancel_immediately",
            type: "boolean",
            required: false,
            description: "Whether to cancel immediately (default: false)",
            location: "body" as const,
          },
        ],
        requestBody: {
          example: `{
  "cancel_immediately": false
}`,
        },
        response: {
          example: `{
  "id": "sub_1234567890",
  "user_id": "usr_1234567890",
  "plan_id": "plan_basic_monthly",
  "status": "active",
  "cancel_at_period_end": true,
  "canceled_at": "2024-01-17T10:30:00Z",
  "current_period_end": "2024-02-01T00:00:00Z"
}`,
        },
      },
    ],
  },
  {
    id: "entitlements",
    title: "Entitlements",
    endpoints: [
      {
        id: "get-user-entitlement",
        title: "Get User Entitlement",
        method: "GET",
        path: "/v1/users/{user_id}/entitlements",
        description:
          "Check if a user has access to a specific entitlement/feature. This is commonly used for feature flagging and access control.",
        authentication: true,
        parameters: [
          {
            name: "user_id",
            type: "string",
            required: true,
            description: "The ID of the user",
            location: "path" as const,
          },
          {
            name: "entitlement_id",
            type: "string",
            required: false,
            description: "Filter by specific entitlement ID",
            location: "query" as const,
          },
        ],
        response: {
          example: `{
  "data": [
    {
      "id": "ent_api_access",
      "name": "API Access",
      "has_access": true,
      "limit": 10000,
      "usage": 2543,
      "reset_at": "2024-02-01T00:00:00Z"
    },
    {
      "id": "ent_advanced_features",
      "name": "Advanced Features",
      "has_access": true,
      "limit": null,
      "usage": null,
      "reset_at": null
    }
  ]
}`,
        },
      },
    ],
  },
  {
    id: "features",
    title: "Features",
    endpoints: [
      {
        id: "list-features",
        title: "List Features",
        method: "GET",
        path: "/v1/features",
        description:
          "Get a list of all available features in your Portcall account. Features are used to define what capabilities are available in different plans.",
        authentication: true,
        parameters: [
          {
            name: "limit",
            type: "number",
            required: false,
            description: "Maximum number of results to return",
            location: "query" as const,
          },
        ],
        response: {
          example: `{
  "data": [
    {
      "id": "feat_api_access",
      "name": "API Access",
      "description": "Access to REST API",
      "type": "metered",
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "feat_premium_support",
      "name": "Premium Support",
      "description": "24/7 priority support",
      "type": "boolean",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}`,
        },
      },
    ],
  },
  {
    id: "plans",
    title: "Plans",
    endpoints: [
      {
        id: "get-plan",
        title: "Get Plan",
        method: "GET",
        path: "/v1/plans/{plan_id}",
        description:
          "Retrieve details about a specific pricing plan, including features, pricing, and billing intervals.",
        authentication: true,
        parameters: [
          {
            name: "plan_id",
            type: "string",
            required: true,
            description: "The ID of the plan to retrieve",
            location: "path" as const,
          },
        ],
        response: {
          example: `{
  "id": "plan_basic_monthly",
  "name": "Basic Plan",
  "description": "Perfect for small teams",
  "price": 2900,
  "currency": "usd",
  "interval": "month",
  "features": [
    {
      "id": "feat_api_access",
      "limit": 10000
    },
    {
      "id": "feat_premium_support",
      "enabled": false
    }
  ],
  "active": true,
  "created_at": "2024-01-01T00:00:00Z"
}`,
        },
      },
    ],
  },
  {
    id: "users",
    title: "Users",
    endpoints: [
      {
        id: "list-users",
        title: "List Users",
        method: "GET",
        path: "/v1/users",
        description:
          "Retrieve a list of users in your Portcall account. Useful for admin dashboards and user management.",
        authentication: true,
        parameters: [
          {
            name: "limit",
            type: "number",
            required: false,
            description:
              "Maximum number of results to return (default: 20, max: 100)",
            location: "query" as const,
          },
          {
            name: "offset",
            type: "number",
            required: false,
            description: "Number of results to skip (for pagination)",
            location: "query" as const,
          },
          {
            name: "email",
            type: "string",
            required: false,
            description: "Filter by email address",
            location: "query" as const,
          },
        ],
        response: {
          example: `{
  "data": [
    {
      "id": "usr_1234567890",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2024-01-01T00:00:00Z",
      "subscription_status": "active"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}`,
        },
      },
    ],
  },
  {
    id: "meter-events",
    title: "Meter Events",
    endpoints: [
      {
        id: "create-meter-event",
        title: "Create Meter Event",
        method: "POST",
        path: "/v1/meter-events",
        description:
          "Record a usage event for metered billing. This is used to track consumption of metered features like API calls, storage, or compute time.",
        authentication: true,
        parameters: [
          {
            name: "user_id",
            type: "string",
            required: true,
            description: "The ID of the user consuming the resource",
            location: "body" as const,
          },
          {
            name: "feature_id",
            type: "string",
            required: true,
            description: "The ID of the metered feature",
            location: "body" as const,
          },
          {
            name: "quantity",
            type: "number",
            required: true,
            description:
              "The amount consumed (e.g., 1 for one API call, 500 for 500MB)",
            location: "body" as const,
          },
          {
            name: "timestamp",
            type: "string",
            required: false,
            description:
              "ISO 8601 timestamp of when the event occurred (defaults to now)",
            location: "body" as const,
          },
          {
            name: "metadata",
            type: "object",
            required: false,
            description: "Additional metadata about the event",
            location: "body" as const,
          },
        ],
        requestBody: {
          example: `{
  "user_id": "usr_1234567890",
  "feature_id": "feat_api_access",
  "quantity": 1,
  "timestamp": "2024-01-17T10:30:00Z",
  "metadata": {
    "endpoint": "/v1/data",
    "method": "GET"
  }
}`,
        },
        response: {
          example: `{
  "id": "evt_1234567890",
  "user_id": "usr_1234567890",
  "feature_id": "feat_api_access",
  "quantity": 1,
  "timestamp": "2024-01-17T10:30:00Z",
  "processed": true
}`,
        },
      },
    ],
  },
];
