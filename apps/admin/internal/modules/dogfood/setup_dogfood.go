package dogfood

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

const (
	DogfoodAccountEmail         = "dogfood@useportcall.com"
	DogfoodLiveAppName          = "Portcall Live"
	DogfoodTestAppName          = "Portcall Test"
	DogfoodFeatureSubscriptions = "max_subscriptions"
	DogfoodFeatureNumberOfUsers = "number_of_users" // Feature for tracking user creation in dashboard
)

type SetupDogfoodRequest struct {
	UpdateK8s bool `json:"update_k8s"`
}

type SetupDogfoodResponse struct {
	Account    AccountInfo `json:"account"`
	LiveApp    AppInfo     `json:"live_app"`
	TestApp    AppInfo     `json:"test_app"`
	LiveSecret string      `json:"live_secret"`
	TestSecret string      `json:"test_secret"`
	Plan       PlanInfo    `json:"plan"`
	Feature    FeatureInfo `json:"feature"`
	K8sUpdated bool        `json:"k8s_updated,omitempty"`
	K8sMessage string      `json:"k8s_message,omitempty"`
}

type AccountInfo struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type AppInfo struct {
	ID       uint   `json:"id"`
	PublicID string `json:"public_id"`
	Name     string `json:"name"`
	IsLive   bool   `json:"is_live"`
}

type PlanInfo struct {
	ID       uint   `json:"id"`
	PublicID string `json:"public_id"`
	Name     string `json:"name"`
}

type FeatureInfo struct {
	ID       uint   `json:"id"`
	PublicID string `json:"public_id"`
}

// SetupDogfood creates or retrieves the dogfood account with live and test apps
// This is used to manage billing/entitlements for dashboard users
func SetupDogfood(c *routerx.Context) {
	response := SetupDogfoodResponse{}

	// Parse optional request body for K8s update flag
	var body SetupDogfoodRequest
	_ = c.ShouldBindJSON(&body) // Ignore error, defaults to false

	// 1. Find or create dogfood account
	account, err := findOrCreateDogfoodAccount(c.DB())
	if err != nil {
		c.ServerError("Failed to setup dogfood account", err)
		return
	}
	response.Account = AccountInfo{ID: account.ID, Email: account.Email}

	// 2. Find or create live app (billing exempt)
	liveApp, err := findOrCreateDogfoodApp(c.DB(), account.ID, DogfoodLiveAppName, true)
	if err != nil {
		c.ServerError("Failed to setup live app", err)
		return
	}
	response.LiveApp = AppInfo{ID: liveApp.ID, PublicID: liveApp.PublicID, Name: liveApp.Name, IsLive: liveApp.IsLive}

	// 3. Find or create test app (billing exempt)
	testApp, err := findOrCreateDogfoodApp(c.DB(), account.ID, DogfoodTestAppName, false)
	if err != nil {
		c.ServerError("Failed to setup test app", err)
		return
	}
	response.TestApp = AppInfo{ID: testApp.ID, PublicID: testApp.PublicID, Name: testApp.Name, IsLive: testApp.IsLive}

	// 4. Create or get secret for live app
	liveSecret, err := findOrCreateSecret(c, liveApp.ID)
	if err != nil {
		c.ServerError("Failed to create live app secret", err)
		return
	}
	response.LiveSecret = liveSecret

	// 5. Create or get secret for test app
	testSecret, err := findOrCreateSecret(c, testApp.ID)
	if err != nil {
		c.ServerError("Failed to create test app secret", err)
		return
	}
	response.TestSecret = testSecret

	// 6. Create the max_subscriptions feature for the live app
	feature, err := findOrCreateFeature(c.DB(), liveApp.ID, DogfoodFeatureSubscriptions)
	if err != nil {
		c.ServerError("Failed to create feature", err)
		return
	}
	response.Feature = FeatureInfo{ID: feature.ID, PublicID: feature.PublicID}

	// 7. Create the number_of_users feature for tracking user creation
	userFeature, err := findOrCreateFeature(c.DB(), liveApp.ID, DogfoodFeatureNumberOfUsers)
	if err != nil {
		log.Printf("Warning: Failed to create number_of_users feature: %v", err)
	}

	// 8. Create a default plan with quota for both features
	plan, err := findOrCreateDogfoodPlan(c.DB(), liveApp.ID, feature.ID, userFeature)
	if err != nil {
		c.ServerError("Failed to create plan", err)
		return
	}
	response.Plan = PlanInfo{ID: plan.ID, PublicID: plan.PublicID, Name: plan.Name}

	// 8b. Create Pro tier plan with more generous limits
	_, err = findOrCreateDogfoodProPlan(c.DB(), liveApp.ID, feature.ID, userFeature)
	if err != nil {
		log.Printf("Warning: Failed to create Pro tier plan: %v", err)
	}

	// 9. Update K8s secrets if requested and secrets are newly created
	if body.UpdateK8s && !isSecretPlaceholder(liveSecret) && !isSecretPlaceholder(testSecret) {
		k8sUpdated, k8sMsg := UpdateK8sSecretsWithNewKeys(liveSecret, testSecret)
		response.K8sUpdated = k8sUpdated
		response.K8sMessage = k8sMsg
	}

	log.Printf("Dogfood setup complete: account=%s, live_app=%s, test_app=%s",
		account.Email, liveApp.PublicID, testApp.PublicID)

	c.OK(response)
}

// isSecretPlaceholder checks if the secret is a placeholder (already exists message)
func isSecretPlaceholder(secret string) bool {
	return secret == "" || secret == "(secret already exists - check database)"
}

func findOrCreateDogfoodAccount(db dbx.IORM) (*models.Account, error) {
	var account models.Account
	err := db.FindFirst(&account, "email = ?", DogfoodAccountEmail)
	if err == nil {
		return &account, nil
	}

	if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}

	// Create new account
	account = models.Account{
		Email:     DogfoodAccountEmail,
		FirstName: "Portcall",
		LastName:  "Dogfood",
	}
	if err := db.Create(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func findOrCreateDogfoodApp(db dbx.IORM, accountID uint, name string, isLive bool) (*models.App, error) {
	var app models.App
	err := db.FindFirst(&app, "account_id = ? AND name = ?", accountID, name)
	if err == nil {
		// Ensure billing exempt is set
		if !app.BillingExempt {
			app.BillingExempt = true
			if err := db.Save(&app); err != nil {
				return nil, err
			}
		}
		return &app, nil
	}

	if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}

	// Create new app
	app = models.App{
		PublicID:      dbx.GenPublicID("app"),
		Name:          name,
		AccountID:     accountID,
		IsLive:        isLive,
		Status:        "active",
		BillingExempt: true, // Dogfood apps are always billing exempt
	}
	if err := db.Create(&app); err != nil {
		return nil, err
	}

	return &app, nil
}

func findOrCreateSecret(c *routerx.Context, appID uint) (string, error) {
	// Check if secret already exists
	var existingSecret models.Secret
	err := c.DB().FindFirst(&existingSecret, "app_id = ? AND disabled_at IS NULL AND key_type = ?", appID, "api_key")
	if err == nil {
		// Secret exists but we can't retrieve the original key
		// Return empty string to indicate secret exists
		return "(secret already exists - check database)", nil
	}

	if !dbx.IsRecordNotFoundError(err) {
		return "", err
	}

	// Generate new API key
	apiKey, err := generateAPIKey(64)
	if err != nil {
		return "", err
	}

	publicID := dbx.GenPublicID("sk")

	hash, err := c.Crypto().Encrypt(apiKey)
	if err != nil {
		return "", err
	}

	secret := models.Secret{
		PublicID: publicID,
		AppID:    appID,
		KeyHash:  hash,
		KeyType:  "api_key",
	}

	if err := c.DB().Create(&secret); err != nil {
		return "", err
	}

	// Return full key (only time it's visible)
	return fmt.Sprintf("%s_%s", publicID, apiKey), nil
}

func findOrCreateFeature(db dbx.IORM, appID uint, featurePublicID string) (*models.Feature, error) {
	var feature models.Feature
	err := db.FindFirst(&feature, "app_id = ? AND public_id = ?", appID, featurePublicID)
	if err == nil {
		return &feature, nil
	}

	if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}

	// Create new feature
	feature = models.Feature{
		PublicID:  featurePublicID,
		AppID:     appID,
		IsMetered: true, // Subscriptions are metered
	}
	if err := db.Create(&feature); err != nil {
		return nil, err
	}

	return &feature, nil
}

func findOrCreateDogfoodPlan(db dbx.IORM, appID uint, subscriptionFeatureID uint, userFeature *models.Feature) (*models.Plan, error) {
	planName := "Dashboard Free Tier"

	var plan models.Plan
	err := db.FindFirst(&plan, "app_id = ? AND name = ?", appID, planName)
	if err == nil {
		// Plan exists, but ensure all plan features are set up
		ensurePlanFeatures(db, &plan, subscriptionFeatureID, userFeature)
		return &plan, nil
	}

	if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}

	// Create new plan
	plan = models.Plan{
		PublicID:      dbx.GenPublicID("plan"),
		AppID:         appID,
		Name:          planName,
		Status:        "published",
		Interval:      "month",
		IntervalCount: 1,
		Currency:      "USD",
		IsFree:        true,
	}
	if err := db.Create(&plan); err != nil {
		return nil, err
	}

	// Create plan item
	planItem := models.PlanItem{
		PublicID:          dbx.GenPublicID("pi"),
		AppID:             appID,
		PlanID:            plan.ID,
		PricingModel:      "fixed",
		Quantity:          1,
		UnitAmount:        0,
		PublicTitle:       "Free Tier",
		PublicDescription: "Up to 10 subscriptions, 100 users",
		PublicUnitLabel:   "subscription",
	}
	if err := db.Create(&planItem); err != nil {
		return nil, err
	}

	// Create plan feature with quota of 10 subscriptions
	planFeature := models.PlanFeature{
		PublicID:   dbx.GenPublicID("pf"),
		AppID:      appID,
		PlanID:     plan.ID,
		PlanItemID: planItem.ID,
		FeatureID:  subscriptionFeatureID,
		Interval:   "month",
		Quota:      10, // Max 10 subscriptions
		Rollover:   0,
	}
	if err := db.Create(&planFeature); err != nil {
		return nil, err
	}

	// Create plan feature with quota of 100 users
	if userFeature != nil {
		userPlanFeature := models.PlanFeature{
			PublicID:   dbx.GenPublicID("pf"),
			AppID:      appID,
			PlanID:     plan.ID,
			PlanItemID: planItem.ID,
			FeatureID:  userFeature.ID,
			Interval:   "month",
			Quota:      100, // Max 100 users
			Rollover:   0,
		}
		if err := db.Create(&userPlanFeature); err != nil {
			log.Printf("Warning: Failed to create number_of_users plan feature: %v", err)
		}
	}

	return &plan, nil
}

// ensurePlanFeatures makes sure all required plan features exist for an existing plan
func ensurePlanFeatures(db dbx.IORM, plan *models.Plan, subscriptionFeatureID uint, userFeature *models.Feature) {
	// Check if user feature plan feature exists
	if userFeature == nil {
		return
	}

	var existingPF models.PlanFeature
	err := db.FindFirst(&existingPF, "plan_id = ? AND feature_id = ?", plan.ID, userFeature.ID)
	if err == nil {
		// Already exists
		return
	}

	if !dbx.IsRecordNotFoundError(err) {
		log.Printf("Warning: Error checking for existing plan feature: %v", err)
		return
	}

	// Get the plan item
	var planItem models.PlanItem
	if err := db.FindFirst(&planItem, "plan_id = ?", plan.ID); err != nil {
		log.Printf("Warning: Could not find plan item for plan %d: %v", plan.ID, err)
		return
	}

	// Create the missing plan feature
	userPlanFeature := models.PlanFeature{
		PublicID:   dbx.GenPublicID("pf"),
		AppID:      plan.AppID,
		PlanID:     plan.ID,
		PlanItemID: planItem.ID,
		FeatureID:  userFeature.ID,
		Interval:   "month",
		Quota:      100, // Max 100 users
		Rollover:   0,
	}
	if err := db.Create(&userPlanFeature); err != nil {
		log.Printf("Warning: Failed to create number_of_users plan feature: %v", err)
	} else {
		log.Printf("Created number_of_users plan feature for plan %s", plan.PublicID)
	}
}

// findOrCreateDogfoodProPlan creates the Pro tier plan with generous limits
func findOrCreateDogfoodProPlan(db dbx.IORM, appID uint, subscriptionFeatureID uint, userFeature *models.Feature) (*models.Plan, error) {
	planName := "Dashboard Pro Tier"

	var plan models.Plan
	err := db.FindFirst(&plan, "app_id = ? AND name = ?", appID, planName)
	if err == nil {
		// Plan exists, ensure features are set up
		ensureProPlanFeatures(db, &plan, subscriptionFeatureID, userFeature)
		return &plan, nil
	}

	if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}

	// Create Pro plan
	plan = models.Plan{
		PublicID:      dbx.GenPublicID("plan"),
		AppID:         appID,
		Name:          planName,
		Status:        "published",
		Interval:      "month",
		IntervalCount: 1,
		Currency:      "USD",
		IsFree:        false,
	}
	if err := db.Create(&plan); err != nil {
		return nil, err
	}

	// Create plan item
	planItem := models.PlanItem{
		PublicID:          dbx.GenPublicID("pi"),
		AppID:             appID,
		PlanID:            plan.ID,
		PricingModel:      "fixed",
		Quantity:          1,
		UnitAmount:        4900, // $49/month
		PublicTitle:       "Pro Tier",
		PublicDescription: "Unlimited subscriptions and users",
		PublicUnitLabel:   "subscription",
	}
	if err := db.Create(&planItem); err != nil {
		return nil, err
	}

	// Create plan feature with unlimited subscriptions (-1)
	planFeature := models.PlanFeature{
		PublicID:   dbx.GenPublicID("pf"),
		AppID:      appID,
		PlanID:     plan.ID,
		PlanItemID: planItem.ID,
		FeatureID:  subscriptionFeatureID,
		Interval:   "month",
		Quota:      -1, // Unlimited
		Rollover:   0,
	}
	if err := db.Create(&planFeature); err != nil {
		return nil, err
	}

	// Create plan feature with unlimited users
	if userFeature != nil {
		userPlanFeature := models.PlanFeature{
			PublicID:   dbx.GenPublicID("pf"),
			AppID:      appID,
			PlanID:     plan.ID,
			PlanItemID: planItem.ID,
			FeatureID:  userFeature.ID,
			Interval:   "month",
			Quota:      -1, // Unlimited
			Rollover:   0,
		}
		if err := db.Create(&userPlanFeature); err != nil {
			log.Printf("Warning: Failed to create Pro plan user feature: %v", err)
		}
	}

	log.Printf("Created Pro tier plan: %s", plan.PublicID)
	return &plan, nil
}

// ensureProPlanFeatures makes sure all required plan features exist for the Pro plan
func ensureProPlanFeatures(db dbx.IORM, plan *models.Plan, subscriptionFeatureID uint, userFeature *models.Feature) {
	if userFeature == nil {
		return
	}

	var existingPF models.PlanFeature
	err := db.FindFirst(&existingPF, "plan_id = ? AND feature_id = ?", plan.ID, userFeature.ID)
	if err == nil {
		return // Already exists
	}

	if !dbx.IsRecordNotFoundError(err) {
		log.Printf("Warning: Error checking for existing Pro plan feature: %v", err)
		return
	}

	var planItem models.PlanItem
	if err := db.FindFirst(&planItem, "plan_id = ?", plan.ID); err != nil {
		log.Printf("Warning: Could not find plan item for Pro plan %d: %v", plan.ID, err)
		return
	}

	userPlanFeature := models.PlanFeature{
		PublicID:   dbx.GenPublicID("pf"),
		AppID:      plan.AppID,
		PlanID:     plan.ID,
		PlanItemID: planItem.ID,
		FeatureID:  userFeature.ID,
		Interval:   "month",
		Quota:      -1, // Unlimited
		Rollover:   0,
	}
	if err := db.Create(&userPlanFeature); err != nil {
		log.Printf("Warning: Failed to create Pro plan user feature: %v", err)
	} else {
		log.Printf("Created user feature for Pro plan %s", plan.PublicID)
	}
}

func generateAPIKey(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	byteLength := (length + 1) / 2
	if length%2 != 0 {
		byteLength++
	}

	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	hexStr := hex.EncodeToString(bytes)
	if len(hexStr) > length {
		hexStr = hexStr[:length]
	}

	return hexStr, nil
}
