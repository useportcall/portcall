package dogfoodx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

const (
	// Feature IDs for dogfood entitlements
	FeatureMaxSubscriptions = "max_subscriptions"
	FeatureNumberOfUsers    = "number_of_users"
)

// ErrEntitlementNotFound is returned when the entitlement doesn't exist in the API
var ErrEntitlementNotFound = fmt.Errorf("entitlement not found")

var (
	config   *Config
	initOnce sync.Once
)

type Config struct {
	APIURL     string
	LiveSecret string
	TestSecret string
}

func init() {
	initOnce.Do(func() {
		config = &Config{
			APIURL:     getEnvOrDefault("DOGFOOD_API_URL", "http://localhost:9080"),
			LiveSecret: os.Getenv("DOGFOOD_LIVE_SECRET"),
			TestSecret: os.Getenv("DOGFOOD_TEST_SECRET"),
		}
		// Log configuration on startup for debugging
		log.Printf("[dogfoodx] Initialized with API URL: %s", config.APIURL)
		log.Printf("[dogfoodx] Live secret configured: %v", config.LiveSecret != "")
		log.Printf("[dogfoodx] Test secret configured: %v", config.TestSecret != "")
	})
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// GetConfig returns the current dogfood configuration
func GetConfig() *Config {
	return config
}

// SetConfig allows overriding the configuration (for testing)
func SetConfig(c *Config) {
	config = c
}

type EntitlementResponse struct {
	ID       string `json:"id"`
	Usage    int64  `json:"usage"`
	Quota    int64  `json:"quota"`
	Interval string `json:"interval"`
}

type QuotaStatus struct {
	FeatureID   string `json:"feature_id"`
	Usage       int64  `json:"usage"`
	Quota       int64  `json:"quota"`
	Remaining   int64  `json:"remaining"`
	IsExceeded  bool   `json:"is_exceeded"`
	IsUnlimited bool   `json:"is_unlimited"`
}

// CheckSubscriptionQuota checks if the app has exceeded its subscription quota
func CheckSubscriptionQuota(db dbx.IORM, appID uint) (*QuotaStatus, error) {
	var app models.App
	if err := db.FindForID(appID, &app); err != nil {
		return nil, err
	}

	log.Printf("[dogfoodx] CheckSubscriptionQuota for app %s (ID: %d, BillingExempt: %v)", app.PublicID, appID, app.BillingExempt)

	// Billing exempt apps have unlimited quota
	if app.BillingExempt {
		log.Printf("[dogfoodx] App %s is billing exempt, returning unlimited", app.PublicID)
		return &QuotaStatus{
			FeatureID:   FeatureMaxSubscriptions,
			Usage:       0,
			Quota:       -1,
			Remaining:   -1,
			IsExceeded:  false,
			IsUnlimited: true,
		}, nil
	}

	// Always use live secret for production dogfood tracking
	// All dashboard apps are tracked in the live dogfood account regardless of their own test/live mode
	apiKey := config.LiveSecret

	if apiKey == "" {
		// No dogfood keys configured, count local database records as fallback
		log.Println("[dogfoodx] Warning: Dogfood API keys not configured, using local database counts")
		return countLocalSubscriptions(db, appID)
	}

	// Map app.PublicID to a user in the dogfood account
	userID := app.PublicID
	log.Printf("[dogfoodx] Fetching entitlement for user %s, feature %s from API", userID, FeatureMaxSubscriptions)

	// Call the Portcall API to check entitlement
	entitlement, err := getEntitlementFromAPI(apiKey, userID, FeatureMaxSubscriptions)
	if err != nil {
		if err == ErrEntitlementNotFound {
			// User not synced to dogfood or no subscription - fall back to local counts
			log.Printf("[dogfoodx] User %s not found in dogfood, using local database counts", userID)
			return countLocalSubscriptions(db, appID)
		}
		log.Printf("[dogfoodx] Warning: Failed to check entitlement from API: %v", err)
		// On API failure, fall back to local counts
		return countLocalSubscriptions(db, appID)
	}

	log.Printf("[dogfoodx] Got entitlement for %s: usage=%d, quota=%d", userID, entitlement.Usage, entitlement.Quota)

	remaining := entitlement.Quota - entitlement.Usage
	if entitlement.Quota < 0 {
		remaining = -1 // Unlimited
	}

	return &QuotaStatus{
		FeatureID:   FeatureMaxSubscriptions,
		Usage:       entitlement.Usage,
		Quota:       entitlement.Quota,
		Remaining:   remaining,
		IsExceeded:  entitlement.Quota >= 0 && entitlement.Usage >= entitlement.Quota,
		IsUnlimited: entitlement.Quota < 0,
	}, nil
}

// countLocalSubscriptions counts subscriptions from the local database
func countLocalSubscriptions(db dbx.IORM, appID uint) (*QuotaStatus, error) {
	var count int64
	if err := db.Count(&count, &models.Subscription{}, "app_id = ?", appID); err != nil {
		log.Printf("[dogfoodx] Warning: Failed to count local subscriptions: %v", err)
		count = 0
	}

	// Default quota of 10 for Free tier when API not configured
	const defaultQuota int64 = 10

	return &QuotaStatus{
		FeatureID:   FeatureMaxSubscriptions,
		Usage:       count,
		Quota:       defaultQuota,
		Remaining:   defaultQuota - count,
		IsExceeded:  count >= defaultQuota,
		IsUnlimited: false,
	}, nil
}

// IncrementSubscriptionUsage increments the subscription count for an app
// Called when a new subscription is created
func IncrementSubscriptionUsage(db dbx.IORM, appID uint) error {
	var app models.App
	if err := db.FindForID(appID, &app); err != nil {
		return err
	}

	// Billing exempt apps don't need to track usage
	if app.BillingExempt {
		log.Printf("Skipping meter event for billing-exempt app %s", app.PublicID)
		return nil
	}

	// Always use live secret for production dogfood tracking
	apiKey := config.LiveSecret

	if apiKey == "" {
		log.Println("Warning: Dogfood API keys not configured, skipping meter event")
		return nil
	}

	log.Printf("Recording subscription meter event for app %s", app.PublicID)
	return recordMeterEvent(apiKey, app.PublicID, FeatureMaxSubscriptions, 1)
}

// IncrementUserUsage increments the number_of_users count for an app
// This is called when a new user is created in the dashboard
//
// ENTITLEMENT GATE: This function tracks user creation for billing purposes.
// When a new user is created, we record a meter event to increment the usage.
// The entitlement system will check if the quota is exceeded before allowing new users.
func IncrementUserUsage(db dbx.IORM, appID uint) error {
	var app models.App
	if err := db.FindForID(appID, &app); err != nil {
		return err
	}

	// Billing exempt apps (dogfood apps) don't need to track usage
	if app.BillingExempt {
		log.Printf("[dogfoodx] Skipping user meter event for billing-exempt app %s", app.PublicID)
		return nil
	}

	// Always use live secret for production dogfood tracking
	apiKey := config.LiveSecret

	if apiKey == "" {
		log.Println("[dogfoodx] Warning: Dogfood API keys not configured, skipping user meter event")
		return nil
	}

	log.Printf("[dogfoodx] Recording user meter event for app %s (feature: %s)", app.PublicID, FeatureNumberOfUsers)
	return recordMeterEvent(apiKey, app.PublicID, FeatureNumberOfUsers, 1)
}

// CheckUserQuota checks if the app has exceeded its user creation quota
//
// ENTITLEMENT GATE: Call this before creating a new user to check if the
// app has remaining quota for user creation.
func CheckUserQuota(db dbx.IORM, appID uint) (*QuotaStatus, error) {
	var app models.App
	if err := db.FindForID(appID, &app); err != nil {
		return nil, err
	}

	// Billing exempt apps have unlimited quota
	if app.BillingExempt {
		return &QuotaStatus{
			FeatureID:   FeatureNumberOfUsers,
			Usage:       0,
			Quota:       -1,
			Remaining:   -1,
			IsExceeded:  false,
			IsUnlimited: true,
		}, nil
	}

	// Always use live secret for production dogfood tracking
	apiKey := config.LiveSecret

	if apiKey == "" {
		// No dogfood keys configured, count local database records as fallback
		log.Println("[dogfoodx] Warning: Dogfood API keys not configured, using local database counts for users")
		return countLocalUsers(db, appID)
	}

	// Map app.PublicID to a user in the dogfood account
	userID := app.PublicID

	// Call the Portcall API to check entitlement
	entitlement, err := getEntitlementFromAPI(apiKey, userID, FeatureNumberOfUsers)
	if err != nil {
		if err == ErrEntitlementNotFound {
			// User not synced to dogfood or no subscription - fall back to local counts
			log.Printf("[dogfoodx] User %s not found in dogfood for user quota, using local database counts", userID)
			return countLocalUsers(db, appID)
		}
		log.Printf("[dogfoodx] Warning: Failed to check user entitlement from API: %v", err)
		// On API failure, fall back to local counts
		return countLocalUsers(db, appID)
	}

	remaining := entitlement.Quota - entitlement.Usage
	if entitlement.Quota < 0 {
		remaining = -1 // Unlimited
	}

	return &QuotaStatus{
		FeatureID:   FeatureNumberOfUsers,
		Usage:       entitlement.Usage,
		Quota:       entitlement.Quota,
		Remaining:   remaining,
		IsExceeded:  entitlement.Quota >= 0 && entitlement.Usage >= entitlement.Quota,
		IsUnlimited: entitlement.Quota < 0,
	}, nil
}

// countLocalUsers counts users from the local database
func countLocalUsers(db dbx.IORM, appID uint) (*QuotaStatus, error) {
	var count int64
	if err := db.Count(&count, &models.User{}, "app_id = ?", appID); err != nil {
		log.Printf("[dogfoodx] Warning: Failed to count local users: %v", err)
		count = 0
	}

	// Default quota of 100 users for Free tier when API not configured
	const defaultQuota int64 = 100

	return &QuotaStatus{
		FeatureID:   FeatureNumberOfUsers,
		Usage:       count,
		Quota:       defaultQuota,
		Remaining:   defaultQuota - count,
		IsExceeded:  count >= defaultQuota,
		IsUnlimited: false,
	}, nil
}

// EnsureUserInDogfood ensures the app exists as a user in the dogfood account
// This should be called when a new app is created in the dashboard
func EnsureUserInDogfood(db dbx.IORM, app *models.App) error {
	if app.BillingExempt {
		return nil
	}

	// Always use live secret for production dogfood tracking
	apiKey := config.LiveSecret

	if apiKey == "" {
		return nil
	}

	// Try to create the user (will fail silently if exists)
	return createUserInDogfood(apiKey, app.PublicID, app.Name)
}

func getEntitlementFromAPI(apiKey, userID, featureID string) (*EntitlementResponse, error) {
	url := fmt.Sprintf("%s/v1/entitlements/%s/%s", config.APIURL, userID, featureID)
	log.Printf("[dogfoodx] Calling API: GET %s", url)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[dogfoodx] API request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[dogfoodx] API response status: %d", resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		// No entitlement found - return error so caller can fall back to local counts
		log.Printf("[dogfoodx] Entitlement not found for %s/%s", userID, featureID)
		return nil, ErrEntitlementNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[dogfoodx] API error: %d - %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var entitlement EntitlementResponse
	if err := json.NewDecoder(resp.Body).Decode(&entitlement); err != nil {
		return nil, err
	}

	log.Printf("[dogfoodx] Parsed entitlement: ID=%s, Usage=%d, Quota=%d", entitlement.ID, entitlement.Usage, entitlement.Quota)
	return &entitlement, nil
}

func recordMeterEvent(apiKey, userID, featureID string, quantity int64) error {
	url := fmt.Sprintf("%s/v1/meter-events", config.APIURL)

	payload := map[string]any{
		"user_id":    userID,
		"feature_id": featureID,
		"usage":      quantity,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("meter event API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Successfully recorded meter event for user %s, feature %s, quantity %d", userID, featureID, quantity)
	return nil
}

func createUserInDogfood(apiKey, userPublicID, appName string) error {
	url := fmt.Sprintf("%s/v1/users", config.APIURL)

	payload := map[string]any{
		"id":    userPublicID,
		"email": fmt.Sprintf("%s@app.portcall.internal", userPublicID),
		"name":  appName,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 201 = created, 409 = already exists (both are fine)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create user API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// SubscribeUserToFreePlan subscribes a user (app) to the free tier plan in dogfood
// This is called after EnsureUserInDogfood to create entitlements for the app
func SubscribeUserToFreePlan(db dbx.IORM, app *models.App, planPublicID string) error {
	if app.BillingExempt {
		return nil
	}

	// Always use live secret for production dogfood tracking
	apiKey := config.LiveSecret

	if apiKey == "" {
		log.Println("[dogfoodx] Warning: Dogfood API keys not configured, skipping subscription")
		return nil
	}

	log.Printf("[dogfoodx] Subscribing app %s to plan %s", app.PublicID, planPublicID)

	// If no plan ID provided, look up the free plan from the API
	if planPublicID == "" {
		freePlanID, err := getFreePlanFromAPI(apiKey)
		if err != nil {
			return fmt.Errorf("failed to get free plan: %w", err)
		}
		planPublicID = freePlanID
		log.Printf("[dogfoodx] Found free plan ID: %s", planPublicID)
	}

	return createSubscriptionInDogfood(apiKey, app.PublicID, planPublicID)
}

func createSubscriptionInDogfood(apiKey, userPublicID, planPublicID string) error {
	url := fmt.Sprintf("%s/v1/subscriptions", config.APIURL)

	payload := map[string]any{
		"user_id": userPublicID,
		"plan_id": planPublicID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 200/201 = created, 409 = already exists with same plan (all fine)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create subscription API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[dogfoodx] Successfully subscribed user %s to plan %s", userPublicID, planPublicID)
	return nil
}

// SyncAppToDogfood ensures an app exists as a user in dogfood and is subscribed to the free tier
// This is called for syncing existing apps that weren't previously registered
// NOTE: Always uses LIVE dogfood API for consistency (all dashboard apps tracked in production)
func SyncAppToDogfood(db dbx.IORM, app *models.App, planPublicID string) error {
	if app.BillingExempt {
		return nil
	}

	// Always use live secret for dogfood sync - we track all apps in the live dogfood account
	apiKey := config.LiveSecret
	if apiKey == "" {
		return fmt.Errorf("dogfood API keys not configured")
	}

	// Step 1: Create/ensure user exists
	if err := createUserInDogfood(apiKey, app.PublicID, app.Name); err != nil {
		return fmt.Errorf("failed to create user in dogfood: %w", err)
	}

	// Step 2: Subscribe to free plan
	if err := createSubscriptionInDogfood(apiKey, app.PublicID, planPublicID); err != nil {
		return fmt.Errorf("failed to subscribe to plan: %w", err)
	}

	// Step 3: Sync current usage (count existing users and set as meter events)
	var userCount int64
	if err := db.Count(&userCount, &models.User{}, "app_id = ?", app.ID); err == nil && userCount > 0 {
		// Record meter events for existing users
		if err := recordMeterEvent(apiKey, app.PublicID, FeatureNumberOfUsers, userCount); err != nil {
			log.Printf("[dogfoodx] Warning: Failed to sync user usage for app %s: %v", app.PublicID, err)
		} else {
			log.Printf("[dogfoodx] Synced %d existing users for app %s", userCount, app.PublicID)
		}
	}

	// Step 4: Sync subscription count
	var subCount int64
	if err := db.Count(&subCount, &models.Subscription{}, "app_id = ?", app.ID); err == nil && subCount > 0 {
		if err := recordMeterEvent(apiKey, app.PublicID, FeatureMaxSubscriptions, subCount); err != nil {
			log.Printf("[dogfoodx] Warning: Failed to sync subscription usage for app %s: %v", app.PublicID, err)
		} else {
			log.Printf("[dogfoodx] Synced %d existing subscriptions for app %s", subCount, app.PublicID)
		}
	}

	return nil
}

type planResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IsFree bool   `json:"is_free"`
}

// getFreePlanFromAPI fetches the free plan ID from the dogfood API
func getFreePlanFromAPI(apiKey string) (string, error) {
	url := fmt.Sprintf("%s/v1/plans", config.APIURL)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("plans API returned status %d: %s", resp.StatusCode, string(body))
	}

	var plans []planResponse
	if err := json.NewDecoder(resp.Body).Decode(&plans); err != nil {
		return "", err
	}

	// Find the first free plan
	for _, p := range plans {
		if p.IsFree {
			return p.ID, nil
		}
	}

	return "", fmt.Errorf("no free plans found")
}

// SubscriptionInfo holds info about the dogfood subscription for an app
type SubscriptionInfo struct {
	PlanName     string `json:"plan_name"`
	PlanPublicID string `json:"plan_public_id"`
	PlanID       string `json:"plan_id"`
	IsFree       bool   `json:"is_free"`
	IsProTier    bool   `json:"is_pro_tier"`
	ProPlanID    string `json:"pro_plan_id,omitempty"`
}

// GetSubscriptionInfo gets the subscription info for an app from dogfood
func GetSubscriptionInfo(db dbx.IORM, appID uint) (*SubscriptionInfo, error) {
	var app models.App
	if err := db.FindForID(appID, &app); err != nil {
		return nil, err
	}

	// Billing exempt apps have unlimited access
	if app.BillingExempt {
		return &SubscriptionInfo{
			PlanName:     "Enterprise (Billing Exempt)",
			PlanPublicID: "",
			IsFree:       false,
			IsProTier:    true,
		}, nil
	}

	apiKey := config.LiveSecret
	if apiKey == "" {
		return &SubscriptionInfo{
			PlanName:  "Free",
			IsFree:    true,
			IsProTier: false,
		}, nil
	}

	userID := app.PublicID

	// Get active subscription from API
	sub, err := getActiveSubscriptionFromAPI(apiKey, userID)
	if err != nil {
		log.Printf("[dogfoodx] Failed to get subscription for %s: %v", userID, err)
		return &SubscriptionInfo{
			PlanName:  "Free",
			IsFree:    true,
			IsProTier: false,
		}, nil
	}

	// Also get the pro plan ID for upgrade flow
	proPlanID, _ := getProPlanFromAPI(apiKey)

	return &SubscriptionInfo{
		PlanName:     sub.PlanName,
		PlanPublicID: sub.PlanID,
		PlanID:       sub.PlanID,
		IsFree:       sub.IsFree,
		IsProTier:    !sub.IsFree,
		ProPlanID:    proPlanID,
	}, nil
}

type subscriptionAPIResponse struct {
	ID       string `json:"id"`
	PlanID   string `json:"plan_id"`
	PlanName string `json:"plan_name"`
	IsFree   bool   `json:"is_free"`
	Status   string `json:"status"`
	Plan     *struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		IsFree bool   `json:"is_free"`
	} `json:"plan"`
}

func getActiveSubscriptionFromAPI(apiKey, userID string) (*subscriptionAPIResponse, error) {
	url := fmt.Sprintf("%s/v1/subscriptions?user_id=%s&status=active", config.APIURL, userID)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("subscriptions API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []subscriptionAPIResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return &subscriptionAPIResponse{
			PlanName: "Free",
			IsFree:   true,
		}, nil
	}

	sub := &result.Data[0]
	// Extract plan info from nested plan object if available
	if sub.Plan != nil {
		sub.PlanID = sub.Plan.ID
		sub.PlanName = sub.Plan.Name
		sub.IsFree = sub.Plan.IsFree
	}

	return sub, nil
}

func getProPlanFromAPI(apiKey string) (string, error) {
	url := fmt.Sprintf("%s/v1/plans", config.APIURL)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("plans API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []planResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Find the first non-free plan (Pro tier)
	for _, p := range result.Data {
		if !p.IsFree {
			return p.ID, nil
		}
	}

	return "", fmt.Errorf("no pro plans found")
}

// CheckoutSessionResponse holds the checkout session info
type CheckoutSessionResponse struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	CheckoutURL string `json:"checkout_url"`
	SessionID   string `json:"session_id"`
}

// CreateCheckoutSession creates a checkout session in dogfood for upgrading to pro
func CreateCheckoutSession(db dbx.IORM, appID uint, cancelURL, redirectURL string) (*CheckoutSessionResponse, error) {
	var app models.App
	if err := db.FindForID(appID, &app); err != nil {
		return nil, err
	}

	if app.BillingExempt {
		return nil, fmt.Errorf("billing exempt apps cannot upgrade")
	}

	apiKey := config.LiveSecret
	if apiKey == "" {
		return nil, fmt.Errorf("dogfood API keys not configured")
	}

	// Ensure user exists in dogfood before creating checkout session
	userID := app.PublicID
	if err := createUserInDogfood(apiKey, userID, app.Name); err != nil {
		return nil, fmt.Errorf("failed to ensure user exists: %w", err)
	}

	// Get the free plan ID
	freePlanID, err := getFreePlanFromAPI(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get free plan: %w", err)
	}

	// Subscribe user to free plan if they don't have a subscription yet
	// This ensures they have entitlements for quota checks
	if err := createSubscriptionInDogfood(apiKey, userID, freePlanID); err != nil {
		// Ignore 409 conflict - user already has subscription
		log.Printf("[dogfoodx] Note: %v (this is fine if user already subscribed)", err)
	}

	// Get the pro plan
	proPlanID, err := getProPlanFromAPI(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get pro plan: %w", err)
	}

	return createCheckoutSessionInAPI(apiKey, userID, proPlanID, cancelURL, redirectURL)
}

func createCheckoutSessionInAPI(apiKey, userID, planID, cancelURL, redirectURL string) (*CheckoutSessionResponse, error) {
	url := fmt.Sprintf("%s/v1/checkout-sessions", config.APIURL)

	payload := map[string]any{
		"user_id":      userID,
		"plan_id":      planID,
		"cancel_url":   cancelURL,
		"redirect_url": redirectURL,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("checkout session API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &CheckoutSessionResponse{
		ID:          result.Data.ID,
		URL:         result.Data.URL,
		CheckoutURL: result.Data.URL,
		SessionID:   result.Data.ID,
	}, nil
}
