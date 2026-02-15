package harness

import "fmt"

// DashHeaders returns admin API key bypass headers for dashboard requests.
func (h *Harness) DashHeaders() map[string]string {
	return map[string]string{
		"X-Admin-API-Key": TestAdminAPIKey,
		"X-Target-App-ID": fmt.Sprintf("%d", h.AppID),
	}
}

// DashPost makes an authenticated POST to the dashboard API.
func (h *Harness) DashPost(path string, body any) *JSONResponse {
	h.T.Helper()
	return DoJSON(h.T, "POST", h.Dashboard.URL+path, body, h.DashHeaders())
}

// DashGet makes an authenticated GET to the dashboard API.
func (h *Harness) DashGet(path string) *JSONResponse {
	h.T.Helper()
	return DoJSON(h.T, "GET", h.Dashboard.URL+path, nil, h.DashHeaders())
}

// CreatePlanViaAPI creates a plan, updates its name/price, and publishes it.
// Returns (planPublicID, planItemPublicID).
func (h *Harness) CreatePlanViaAPI(name string, cents int64) (string, string) {
	h.T.Helper()
	base := DashPath(h.AppPublicID, "plans")

	data := h.DashPost(base, nil).MustOK(h.T, "create plan")
	planID := MustString(h.T, data, "id")
	items := GetSlice(data, "items")
	if len(items) == 0 {
		h.T.Fatal("create plan returned no items")
	}
	itemID := MustString(h.T, items[0].(map[string]any), "id")

	h.DashPost(base+"/"+planID, map[string]any{
		"name": name, "currency": "usd",
	}).MustOK(h.T, "update plan name")

	h.DashPost(DashPath(h.AppPublicID, "plan-items/"+itemID), map[string]any{
		"unit_amount": cents,
	}).MustOK(h.T, "update plan item price")

	h.DashPost(base+"/"+planID+"/publish", nil).MustOK(h.T, "publish plan")
	return planID, itemID
}

// CreateSecretViaAPI creates an API secret and stores the raw key on h.APIKey.
func (h *Harness) CreateSecretViaAPI() string {
	h.T.Helper()
	data := h.DashPost(DashPath(h.AppPublicID, "secrets"), nil).
		MustOK(h.T, "create secret")
	key := MustString(h.T, data, "key")
	h.APIKey = key
	return key
}

// CreateUserViaAPI creates a user via the public API.
func (h *Harness) CreateUserViaAPI(name, email string) string {
	h.T.Helper()
	data := h.APIPost("/v1/users", map[string]any{
		"name": name, "email": email,
	}).MustOK(h.T, "create user")
	return MustString(h.T, data, "id")
}

// CreateFeatureViaAPI creates a feature via the public API.
func (h *Harness) CreateFeatureViaAPI(featureID string, isMetered bool) string {
	h.T.Helper()
	data := h.APIPost("/v1/features", map[string]any{
		"feature_id": featureID,
		"is_metered": isMetered,
	}).MustOK(h.T, "create feature")
	return MustString(h.T, data, "id")
}

// CreatePlanFeatureViaAPI attaches a feature to a plan via the public API.
func (h *Harness) CreatePlanFeatureViaAPI(featureID, planID string, quota int64) string {
	h.T.Helper()
	data := h.APIPost("/v1/plan-features", map[string]any{
		"feature_id": featureID,
		"plan_id":    planID,
		"quota":      quota,
	}).MustOK(h.T, "create plan feature")
	return MustString(h.T, data, "id")
}

// ListEntitlementsViaAPI lists entitlements for a user via the public API.
func (h *Harness) ListEntitlementsViaAPI(userID string) []map[string]any {
	h.T.Helper()
	return h.APIGet("/v1/entitlements?user_id=" + userID).
		MustOKList(h.T, "list entitlements")
}

// GetEntitlementViaAPI gets a single entitlement by user + feature ID.
func (h *Harness) GetEntitlementViaAPI(userID, featureID string) map[string]any {
	h.T.Helper()
	return h.APIGet("/v1/entitlements/" + userID + "/" + featureID).
		MustOK(h.T, "get entitlement")
}

// RecordMeterEventViaAPI records a usage event via the public API.
func (h *Harness) RecordMeterEventViaAPI(userID, featureID string, usage int) *JSONResponse {
	h.T.Helper()
	return h.APIPost("/v1/meter-events", map[string]any{
		"user_id":    userID,
		"feature_id": featureID,
		"usage":      usage,
	})
}
