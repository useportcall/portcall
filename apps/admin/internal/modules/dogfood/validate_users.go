package dogfood

import (
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type ValidateUsersRequest struct {
	DryRun bool `json:"dry_run"`
}

type UserValidationResult struct {
	UserPublicID string `json:"user_public_id"`
	UserEmail    string `json:"user_email"`
	UserName     string `json:"user_name"`
	AppPublicID  string `json:"app_public_id,omitempty"`
	AppName      string `json:"app_name,omitempty"`
	AppIsLive    bool   `json:"app_is_live,omitempty"`
	AccountEmail string `json:"account_email,omitempty"`
	Status       string `json:"status"`
	Action       string `json:"action,omitempty"`
	Issue        string `json:"issue,omitempty"`
	Fixed        bool   `json:"fixed,omitempty"`
	Error        string `json:"error,omitempty"`
}

type ValidateUsersResponse struct {
	LiveAppID        uint                   `json:"live_app_id"`
	LiveAppName      string                 `json:"live_app_name"`
	TestAppID        uint                   `json:"test_app_id"`
	TestAppName      string                 `json:"test_app_name"`
	TotalUsers       int                    `json:"total_users"`
	Valid            int                    `json:"valid"`
	MissingApp       int                    `json:"missing_app"`
	WrongEnvironment int                    `json:"wrong_environment"`
	EmailMismatch    int                    `json:"email_mismatch"`
	Fixed            int                    `json:"fixed"`
	Failed           int                    `json:"failed"`
	DryRun           bool                   `json:"dry_run"`
	Results          []UserValidationResult `json:"results"`
}

// ValidateUsers validates that dogfood users correspond correctly to apps
// - Live apps should be users in the dogfood live app
// - Test apps should be users in the dogfood test app
// - User email should match the app's account email
func ValidateUsers(c *routerx.Context) {
	var body ValidateUsersRequest
	_ = c.ShouldBindJSON(&body)

	// Find the dogfood account
	var dogfoodAccount models.Account
	if err := c.DB().FindFirst(&dogfoodAccount, "email = ?", DogfoodAccountEmail); err != nil {
		c.ServerError("Dogfood account not found. Run /api/dogfood/setup first.", err)
		return
	}

	// Find the live and test apps
	var liveApp, testApp models.App
	if err := c.DB().FindFirst(&liveApp, "account_id = ? AND name = ?", dogfoodAccount.ID, DogfoodLiveAppName); err != nil {
		c.ServerError("Dogfood live app not found. Run /api/dogfood/setup first.", err)
		return
	}
	if err := c.DB().FindFirst(&testApp, "account_id = ? AND name = ?", dogfoodAccount.ID, DogfoodTestAppName); err != nil {
		c.ServerError("Dogfood test app not found. Run /api/dogfood/setup first.", err)
		return
	}

	response := ValidateUsersResponse{
		LiveAppID:   liveApp.ID,
		LiveAppName: liveApp.Name,
		TestAppID:   testApp.ID,
		TestAppName: testApp.Name,
		DryRun:      body.DryRun,
		Results:     make([]UserValidationResult, 0),
	}

	// Get all users from both dogfood apps
	var liveUsers, testUsers []models.User
	if err := c.DB().List(&liveUsers, "app_id = ?", liveApp.ID); err != nil {
		c.ServerError("Failed to list live app users", err)
		return
	}
	if err := c.DB().List(&testUsers, "app_id = ?", testApp.ID); err != nil {
		c.ServerError("Failed to list test app users", err)
		return
	}

	response.TotalUsers = len(liveUsers) + len(testUsers)

	// Build a map of all non-dogfood apps by their public_id
	// Each "user" in the dogfood app represents an app using Portcall
	var allApps []models.App
	if err := c.DB().List(&allApps, "account_id != ?", dogfoodAccount.ID); err != nil {
		c.ServerError("Failed to list apps", err)
		return
	}

	// Create maps for quick lookup
	appByPublicID := make(map[string]*models.App)
	appByAccountEmail := make(map[string]*models.App)
	for i := range allApps {
		appByPublicID[allApps[i].PublicID] = &allApps[i]
	}

	// Get all accounts for apps (for email lookup)
	accountIDs := make([]uint, 0)
	for _, app := range allApps {
		accountIDs = append(accountIDs, app.AccountID)
	}
	var accounts []models.Account
	if len(accountIDs) > 0 {
		c.DB().List(&accounts, "id IN ?", accountIDs)
	}
	accountByID := make(map[uint]*models.Account)
	for i := range accounts {
		accountByID[accounts[i].ID] = &accounts[i]
	}

	// Build map of apps by their account email for orphan linking
	for publicID, app := range appByPublicID {
		if account := accountByID[app.AccountID]; account != nil {
			appByAccountEmail[account.Email] = appByPublicID[publicID]
		}
	}

	// Deduplicate users first (remove duplicates with same public_id in same app)
	liveUsers = deduplicateUsers(c, liveUsers, body.DryRun)
	testUsers = deduplicateUsers(c, testUsers, body.DryRun)

	// Process live app users
	processUsers(c, &response, liveUsers, &liveApp, &testApp, appByPublicID, appByAccountEmail, accountByID, true, body.DryRun)

	// Process test app users
	processUsers(c, &response, testUsers, &liveApp, &testApp, appByPublicID, appByAccountEmail, accountByID, false, body.DryRun)

	// After validation, ensure all existing apps have users in the dogfood app
	if !body.DryRun {
		createMissingUsers(c, &response, &liveApp, &testApp, allApps, accountByID)
	}

	c.OK(response)
}

func processUsers(
	c *routerx.Context,
	response *ValidateUsersResponse,
	users []models.User,
	liveApp *models.App,
	testApp *models.App,
	appByPublicID map[string]*models.App,
	appByAccountEmail map[string]*models.App,
	accountByID map[uint]*models.Account,
	isLiveContext bool, // true if processing users from the live dogfood app
	dryRun bool,
) {
	for _, user := range users {
		result := UserValidationResult{
			UserPublicID: user.PublicID,
			UserEmail:    user.Email,
			UserName:     user.Name,
		}

		// The user's email or public_id might reference an app
		// Users in dogfood represent apps, so we try to find the corresponding app
		// The user.Email typically is the app's account email or formatted as <app_id>@app.portcall.internal
		app := findAppForUser(&user, appByPublicID, appByAccountEmail)

		if app == nil {
			// Could not find a corresponding app for this dogfood user - orphan
			result.Status = "orphan"
			result.Issue = "No corresponding app found for this dogfood user"
			result.Action = "delete"
			response.MissingApp++

			if !dryRun {
				// Delete orphan users
				if err := c.DB().Delete(&user, "id = ?", user.ID); err != nil {
					result.Error = fmt.Sprintf("Failed to delete orphan user: %v", err)
					response.Failed++
				} else {
					result.Fixed = true
					response.Fixed++
					log.Printf("[dogfood/validate] Deleted orphan user %s (%s)", user.PublicID, user.Email)
				}
			}

			response.Results = append(response.Results, result)
			continue
		}

		result.AppPublicID = app.PublicID
		result.AppName = app.Name
		result.AppIsLive = app.IsLive

		// Get the account email for this app
		account := accountByID[app.AccountID]
		if account != nil {
			result.AccountEmail = account.Email
		}

		// Check 1: Does the user email match the app's account email?
		// This is now the FIRST check - user email must match account email
		emailMismatch := false
		if account != nil && user.Email != account.Email {
			emailMismatch = true
			result.Issue = fmt.Sprintf("User email (%s) doesn't match account email (%s)", user.Email, account.Email)
		}

		if emailMismatch {
			result.Status = "email_mismatch"
			result.Action = "update_email"
			response.EmailMismatch++

			if !dryRun && account != nil {
				if err := updateUserEmail(c.DB(), &user, account.Email); err != nil {
					result.Error = err.Error()
					response.Failed++
				} else {
					result.Fixed = true
					result.UserEmail = account.Email
					response.Fixed++
					log.Printf("[dogfood/validate] Updated user %s email from %s to %s",
						user.PublicID, user.Email, account.Email)
				}
			}

			response.Results = append(response.Results, result)
			continue
		}

		// Check 2: Is the app in the correct dogfood environment?
		// Live apps should be users in the live dogfood app
		// Test apps should be users in the test dogfood app
		wrongEnvironment := false
		if app.IsLive && !isLiveContext {
			wrongEnvironment = true
			result.Issue = "Live app found in test dogfood app"
		} else if !app.IsLive && isLiveContext {
			wrongEnvironment = true
			result.Issue = "Test app found in live dogfood app"
		}

		if wrongEnvironment {
			result.Status = "wrong_environment"
			result.Action = "move_user"
			response.WrongEnvironment++

			if !dryRun {
				// Move user to the correct dogfood app
				targetAppID := liveApp.ID
				if !app.IsLive {
					targetAppID = testApp.ID
				}

				if err := moveUserToApp(c.DB(), &user, targetAppID); err != nil {
					result.Error = err.Error()
					response.Failed++
				} else {
					result.Fixed = true
					response.Fixed++
					log.Printf("[dogfood/validate] Moved user %s to %s app", user.PublicID,
						map[bool]string{true: "live", false: "test"}[app.IsLive])
				}
			}

			response.Results = append(response.Results, result)
			continue
		}

		// All checks passed
		result.Status = "valid"
		response.Valid++
		// Don't add valid users to results to keep the response smaller
	}
}

// findAppForUser tries to find the corresponding app for a dogfood user
func findAppForUser(user *models.User, appByPublicID map[string]*models.App, appByAccountEmail map[string]*models.App) *models.App {
	// Try to match by user's public_id (if it looks like an app public_id)
	if app, ok := appByPublicID[user.PublicID]; ok {
		return app
	}

	// Try to match by extracting app_id from internal email format
	// e.g., "app_xxxxx@app.portcall.internal" -> "app_xxxxx"
	if isInternalEmail(user.Email) {
		appID := extractAppIDFromEmail(user.Email)
		if app, ok := appByPublicID[appID]; ok {
			return app
		}
	}

	// Try to match by user's email to account email (for orphan linking)
	if app, ok := appByAccountEmail[user.Email]; ok {
		return app
	}

	// Try to match by user's name if it contains an app public_id
	if app, ok := appByPublicID[user.Name]; ok {
		return app
	}

	return nil
}

// isInternalEmail checks if the email is an internal portcall email
func isInternalEmail(email string) bool {
	return len(email) > 22 && email[len(email)-22:] == "@app.portcall.internal"
}

// extractAppIDFromEmail extracts the app_id from an internal email
func extractAppIDFromEmail(email string) string {
	if !isInternalEmail(email) {
		return ""
	}
	// Remove the "@app.portcall.internal" suffix
	return email[:len(email)-22]
}

// moveUserToApp moves a user to a different dogfood app
func moveUserToApp(db dbx.IORM, user *models.User, targetAppID uint) error {
	user.AppID = targetAppID
	return db.Save(user)
}

// updateUserEmail updates a user's email
func updateUserEmail(db dbx.IORM, user *models.User, newEmail string) error {
	user.Email = newEmail
	return db.Save(user)
}

// deduplicateUsers removes duplicate users with the same public_id
// Keeps the oldest user (by ID) and deletes the rest
func deduplicateUsers(c *routerx.Context, users []models.User, dryRun bool) []models.User {
	// Group users by public_id
	usersByPublicID := make(map[string][]models.User)
	for _, user := range users {
		usersByPublicID[user.PublicID] = append(usersByPublicID[user.PublicID], user)
	}

	// Find and remove duplicates
	deduplicatedUsers := make([]models.User, 0)
	for _, userGroup := range usersByPublicID {
		if len(userGroup) == 1 {
			// No duplicates
			deduplicatedUsers = append(deduplicatedUsers, userGroup[0])
			continue
		}

		// Found duplicates - keep the one with the lowest ID (oldest)
		var keepUser *models.User
		duplicates := make([]models.User, 0)

		for i := range userGroup {
			if keepUser == nil || userGroup[i].ID < keepUser.ID {
				if keepUser != nil {
					duplicates = append(duplicates, *keepUser)
				}
				keepUser = &userGroup[i]
			} else {
				duplicates = append(duplicates, userGroup[i])
			}
		}

		// Delete duplicates
		for _, dup := range duplicates {
			if !dryRun {
				if err := c.DB().Delete(&dup, "id = ?", dup.ID); err != nil {
					log.Printf("[dogfood/validate] Failed to delete duplicate user %s (ID: %d): %v", dup.PublicID, dup.ID, err)
				} else {
					log.Printf("[dogfood/validate] Deleted duplicate user %s (ID: %d, Email: %s), kept user ID: %d",
						dup.PublicID, dup.ID, dup.Email, keepUser.ID)
				}
			} else {
				log.Printf("[dogfood/validate] [DRY RUN] Would delete duplicate user %s (ID: %d, Email: %s), keeping user ID: %d",
					dup.PublicID, dup.ID, dup.Email, keepUser.ID)
			}
		}

		// Add the kept user to the deduplicated list
		if keepUser != nil {
			deduplicatedUsers = append(deduplicatedUsers, *keepUser)
		}
	}

	return deduplicatedUsers
}

// createMissingUsers ensures all apps have corresponding users in the dogfood app
func createMissingUsers(
	c *routerx.Context,
	response *ValidateUsersResponse,
	liveApp *models.App,
	testApp *models.App,
	allApps []models.App,
	accountByID map[uint]*models.Account,
) {
	// Get existing users from both dogfood apps
	var existingLiveUsers, existingTestUsers []models.User
	c.DB().List(&existingLiveUsers, "app_id = ?", liveApp.ID)
	c.DB().List(&existingTestUsers, "app_id = ?", testApp.ID)

	// Build sets of app public_ids that already have users
	liveUserAppIDs := make(map[string]bool)
	testUserAppIDs := make(map[string]bool)

	for _, user := range existingLiveUsers {
		liveUserAppIDs[user.PublicID] = true
	}
	for _, user := range existingTestUsers {
		testUserAppIDs[user.PublicID] = true
	}

	// Process live apps first, then test apps
	liveApps := make([]models.App, 0)
	testApps := make([]models.App, 0)

	for _, app := range allApps {
		if app.IsLive {
			liveApps = append(liveApps, app)
		} else {
			testApps = append(testApps, app)
		}
	}

	// Create users for live apps
	for _, app := range liveApps {
		if liveUserAppIDs[app.PublicID] {
			continue // User already exists
		}

		account := accountByID[app.AccountID]
		if account == nil {
			log.Printf("[dogfood/validate] Skipping app %s - no account found", app.PublicID)
			continue
		}

		// Create user in live dogfood app
		user := models.User{
			PublicID: app.PublicID,
			AppID:    liveApp.ID,
			Name:     app.Name,
			Email:    account.Email,
		}

		if err := c.DB().Create(&user); err != nil {
			log.Printf("[dogfood/validate] Failed to create user for app %s: %v", app.PublicID, err)
		} else {
			log.Printf("[dogfood/validate] Created user for live app %s (%s)", app.PublicID, account.Email)
		}
	}

	// Create users for test apps
	for _, app := range testApps {
		if testUserAppIDs[app.PublicID] {
			continue // User already exists
		}

		account := accountByID[app.AccountID]
		if account == nil {
			log.Printf("[dogfood/validate] Skipping app %s - no account found", app.PublicID)
			continue
		}

		// Create user in test dogfood app
		user := models.User{
			PublicID: app.PublicID,
			AppID:    testApp.ID,
			Name:     app.Name,
			Email:    account.Email,
		}

		if err := c.DB().Create(&user); err != nil {
			log.Printf("[dogfood/validate] Failed to create user for app %s: %v", app.PublicID, err)
		} else {
			log.Printf("[dogfood/validate] Created user for test app %s (%s)", app.PublicID, account.Email)
		}
	}
}
