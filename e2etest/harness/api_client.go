package harness

import "net/url"

// apiHeaders returns the API key auth header for public API requests.
func (h *Harness) apiHeaders() map[string]string {
	if h.APIKey == "" {
		h.T.Fatal("APIKey not set — call CreateSecretViaAPI first")
	}
	return map[string]string{"x-api-key": h.APIKey}
}

// APIPost makes an authenticated POST to the public API.
func (h *Harness) APIPost(path string, body any) *JSONResponse {
	h.T.Helper()
	return DoJSON(h.T, "POST", h.API.URL+path, body, h.apiHeaders())
}

// APIGet makes an authenticated GET to the public API.
func (h *Harness) APIGet(path string) *JSONResponse {
	h.T.Helper()
	return DoJSON(h.T, "GET", h.API.URL+path, nil, h.apiHeaders())
}

// CheckoutGet makes a GET to the checkout API with session token auth.
func (h *Harness) CheckoutGet(path, sessionToken string) *JSONResponse {
	h.T.Helper()
	hdrs := map[string]string{"X-Checkout-Session-Token": sessionToken}
	return DoJSON(h.T, "GET", h.Checkout.URL+path, nil, hdrs)
}

// CreateCheckoutSessionViaAPI creates a checkout session via the public API.
// Returns (response data, session token extracted from the checkout URL).
func (h *Harness) CreateCheckoutSessionViaAPI(planID, userID string) (map[string]any, string) {
	h.T.Helper()
	data := h.APIPost("/v1/checkout-sessions", map[string]any{
		"plan_id":      planID,
		"user_id":      userID,
		"cancel_url":   "https://e2e.test/cancel",
		"redirect_url": "https://e2e.test/success",
	}).MustOK(h.T, "create checkout session")

	csURL := GetString(data, "url")
	parsed, err := url.Parse(csURL)
	if err != nil {
		h.T.Fatalf("parse checkout URL %q: %v", csURL, err)
	}
	token := parsed.Query().Get("st")
	if token == "" {
		h.T.Fatalf("checkout URL missing st param: %s", csURL)
	}
	return data, token
}

// SetBillingAddressViaAPI sets a billing address on a user via the public API.
func (h *Harness) SetBillingAddressViaAPI(userID string) {
	h.T.Helper()
	h.APIPost("/v1/users/"+userID+"/billing-address", map[string]any{
		"line1": "1 Test Blvd", "city": "San Francisco",
		"postal_code": "94105", "country": "US",
	}).MustOK(h.T, "set billing address")
}
