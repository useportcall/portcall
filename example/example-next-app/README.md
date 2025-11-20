# Portcall Next.js Example App

## API Endpoints Called

```bash
# create checkout session
POST https://api.portcall.com/v1/checkout-sessions`

# retrieve user subscriptions
GET https://api.portcall.com/v1/subscriptions?user_id={{user_id}}

# update subscription
POST https://api.portcall.com/v1/subscriptions/{{subscription_id}}

# cancel subscription
POST https://api.portcall.com/v1/subscriptions/{{subscription_id}}/cancel

# retrieve user entitlement
GET https://api.portcall.com/v1/users/{{user_id}}/entitlements?entitlement_id={{entitlement_id}}
```
