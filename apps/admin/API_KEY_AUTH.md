# Admin API Key Authentication

The Admin API supports two authentication methods:

1. **Keycloak JWT** (for UI access)
2. **API Key** (for programmatic access)

## Using API Key Authentication

### Setup

The API key is configured via the `ADMIN_API_KEY` environment variable.

**Default API Key:**

```
admin_e5f7a9b2c4d6e8f0a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7
```

### Making API Requests

Include the API key in the `X-Admin-API-Key` header:

```bash
curl "http://localhost:8081/api/apps" \
  -H "X-Admin-API-Key: admin_e5f7a9b2c4d6e8f0a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7"
```

### Examples

#### Get Dogfood Status

```bash
curl "http://localhost:8081/api/dogfood/status" \
  -H "X-Admin-API-Key: admin_e5f7a9b2c4d6e8f0a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7"
```

#### Setup Dogfood Account

```bash
curl -X POST "http://localhost:8081/api/dogfood/setup" \
  -H "X-Admin-API-Key: admin_e5f7a9b2c4d6e8f0a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7"
```

#### List Apps

```bash
curl "http://localhost:8081/api/apps" \
  -H "X-Admin-API-Key: admin_e5f7a9b2c4d6e8f0a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7"
```

#### Get App Details

```bash
curl "http://localhost:8081/api/apps/1" \
  -H "X-Admin-API-Key: admin_e5f7a9b2c4d6e8f0a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7"
```

## Security Notes

⚠️ **Important Security Considerations:**

- **Change the default API key** in production environments
- Store the API key securely (use secrets management)
- Rotate the API key periodically
- Never commit the API key to version control
- Use HTTPS in production to protect the key in transit
- API keys must start with `admin_` prefix

## Generating a New API Key

To generate a new secure API key:

```bash
# Generate a random 64-character hex string
openssl rand -hex 32 | awk '{print "admin_" $0}'
```

Then update the `ADMIN_API_KEY` environment variable in:

- `/apps/admin/.env` (local development)
- `/docker-compose/docker-compose.admin.yml` (Docker)
- Production deployment configuration

## Error Responses

### Invalid API Key

```json
{
  "access": "unauthorized",
  "error": "Invalid API key"
}
```

### Missing API Key (falls back to Keycloak)

```json
{
  "access": "unauthorized"
}
```

### API Key Not Configured

```json
{
  "access": "unauthorized",
  "error": "API key authentication not configured"
}
```

## Implementation Details

The authentication middleware checks for the `X-Admin-API-Key` header first:

- If present and valid → authenticated as `admin@portcall.internal`
- If present but invalid → returns 401 Unauthorized
- If not present → falls back to Keycloak JWT validation

This allows both the UI (via Keycloak) and programmatic access (via API key) to work simultaneously.
