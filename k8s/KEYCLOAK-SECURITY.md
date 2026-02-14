# Keycloak Security Configuration

## 🔒 Security Measures Implemented

### 1. Strong Admin Password

- **Default password changed**: The weak default `admin_change_me_123` has been replaced with a strong 32-character random password
- **Stored securely**: Password is stored in Kubernetes secret `portcall-secrets`
- **Script**: [scripts/secure-keycloak-admin.sh](../scripts/secure-keycloak-admin.sh)

### 2. Network Policy Lockdown

- **File**: [k8s/portcall-chart/templates/keycloak-network-policy.yaml](../k8s/portcall-chart/templates/keycloak-network-policy.yaml)
- **Blocks**: All external access to admin API endpoints
- **Allows**:
  - Health checks from within cluster
  - User-facing endpoints (login, token exchange) via ingress
  - Internal pod-to-pod communication

### 3. No Public Admin Access

- **Admin console is NOT publicly accessible**
- **Configuration**: `keycloak.adminAccess.enabled: false` in values.yaml
- **Emergency access**: Can be temporarily enabled with IP whitelisting (see below)

---

## 🔑 Retrieving Admin Credentials

### Get Admin Username

```bash
kubectl get secret portcall-secrets -n portcall -o jsonpath='{.data.KC_BOOTSTRAP_ADMIN_USERNAME}' | base64 -d
echo
```

### Get Admin Password

```bash
kubectl get secret portcall-secrets -n portcall -o jsonpath='{.data.KC_BOOTSTRAP_ADMIN_PASSWORD}' | base64 -d
echo
```

### Get Both (One Command)

```bash
echo "Username: $(kubectl get secret portcall-secrets -n portcall -o jsonpath='{.data.KC_BOOTSTRAP_ADMIN_USERNAME}' | base64 -d)"
echo "Password: $(kubectl get secret portcall-secrets -n portcall -o jsonpath='{.data.KC_BOOTSTRAP_ADMIN_PASSWORD}' | base64 -d)"
```

---

## 🖥️ Accessing Keycloak Admin Console

### Method 1: kubectl port-forward (Recommended)

**This is the SECURE way to access the admin console.**

```bash
# Forward Keycloak admin console to local port 8080
kubectl port-forward -n portcall deployment/keycloak 8080:8080

# In another terminal, open your browser or use curl
open http://localhost:8080/admin

# Login with credentials from secret (see above)
```

**Advantages:**

- ✅ No public exposure
- ✅ Traffic encrypted within kubectl tunnel
- ✅ Automatically closes when you stop the command
- ✅ No firewall rules needed

---

### Method 2: Temporary Public Access (Emergency Only)

**⚠️ ONLY use this if you MUST access admin console from a remote location**

1. **Get your current public IP:**

```bash
curl -s https://api.ipify.org
```

2. **Update values.yaml:**

```yaml
keycloak:
  adminAccess:
    enabled: true # Temporarily enable
    host: "admin.useportcall.com" # Separate subdomain for admin
    allowedIPs:
      - "YOUR.PUBLIC.IP.HERE/32" # Your IP only
```

3. **Deploy changes:**

```bash
cd k8s
helm upgrade portcall ./portcall-chart \
  --namespace portcall \
  --values ./deploy/digitalocean/values.yaml
```

4. **Access admin console:**

```
https://admin.useportcall.com/admin
```

5. **IMMEDIATELY disable after use:**

```yaml
keycloak:
  adminAccess:
    enabled: false # Disable again
```

Then re-deploy:

```bash
helm upgrade portcall ./portcall-chart \
  --namespace portcall \
  --values ./deploy/digitalocean/values.yaml
```

---

## 🛡️ Security Best Practices

### DO ✅

- Use `kubectl port-forward` for admin access
- Store admin password in password manager
- Rotate admin password regularly
- Use service accounts with limited permissions for automation
- Review Keycloak audit logs regularly

### DON'T ❌

- Leave `adminAccess.enabled: true` permanently
- Use the admin account for application integrations
- Share admin credentials
- Expose admin console without IP whitelisting
- Use weak or default passwords

---

## 🔄 Rotating Admin Password

### Manual Rotation

```bash
# Run the secure admin script again
./scripts/secure-keycloak-admin.sh

# This will:
# 1. Generate a new random password
# 2. Update Keycloak admin user
# 3. Update Kubernetes secret
# 4. Restart Keycloak deployment
```

### Scheduled Rotation (Recommended)

```bash
# Create a CronJob to rotate password monthly
kubectl create cronjob keycloak-password-rotation \
  --schedule="0 0 1 * *" \
  --image=curlimages/curl:latest \
  --namespace=portcall \
  -- /scripts/secure-keycloak-admin.sh
```

---

## 🚨 Incident Response

### If Admin Credentials Are Compromised

1. **Immediately rotate password:**

```bash
./scripts/secure-keycloak-admin.sh
```

2. **Check audit logs for unauthorized access:**

```bash
kubectl logs -n portcall deployment/keycloak | grep -i admin
```

3. **Review recent user/client changes:**

```bash
# Port-forward to admin console
kubectl port-forward -n portcall deployment/keycloak 8080:8080

# Navigate to: Admin Console → Events → Login Events
# Check for suspicious admin logins
```

4. **Disable admin access if enabled:**

```yaml
keycloak:
  adminAccess:
    enabled: false
```

5. **Review and rotate all client secrets:**
   - Navigate to Clients → Credentials
   - Regenerate secrets for all OAuth clients

---

## 📊 Monitoring Admin Access

### View Recent Admin API Calls

```bash
kubectl logs -n portcall deployment/keycloak --tail=100 | grep "/admin/"
```

### Monitor Failed Login Attempts

```bash
kubectl logs -n portcall deployment/keycloak | grep -i "failed.*admin"
```

### Set Up Alerts

Configure Prometheus alerts for:

- Failed admin login attempts
- Admin console access from unexpected IPs
- Unusual admin API activity

---

## 🧪 Testing Password Reset (Users)

The password reset feature is now configured for end users:

```bash
# Test password reset for a user
./scripts/test-password-reset.sh

# Check Postmark Activity dashboard
open https://account.postmarkapp.com/servers/YOUR-SERVER-ID/streams/outbound/activity
```

**Configuration:**

- SMTP: Postmark (smtp.postmarkapp.com:587)
- From: notifications@mail.useportcall.com
- Enabled in realm: ✅ "Forgot Password" feature enabled

---

## 📝 Summary

| Feature                 | Status        | Details                             |
| ----------------------- | ------------- | ----------------------------------- |
| **Admin Password**      | 🔒 Secured    | Strong random 32-char password      |
| **Public Admin Access** | ❌ Disabled   | No external access to admin console |
| **Network Policy**      | ✅ Active     | Only allows necessary traffic       |
| **Port Forward Access** | ✅ Available  | Secure local access method          |
| **Emergency Access**    | ⚠️ Available  | Can be enabled with IP whitelist    |
| **Password Reset**      | ✅ Configured | Users can reset via email           |

**Current Security Posture**: 🟢 **SECURE**

Third parties **CANNOT** access or modify Keycloak configuration. Admin access is only available via secure kubectl port-forward or with explicit IP whitelisting.
