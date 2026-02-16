package dogfood

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type FixUsersRequest struct {
	DryRun bool `json:"dry_run"`
}

type UserFixResult struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Action     string `json:"action"`
	OldEmail   string `json:"old_email,omitempty"`
	NewEmail   string `json:"new_email,omitempty"`
	Subscribed bool   `json:"subscribed,omitempty"`
	PlanName   string `json:"plan_name,omitempty"`
	Error      string `json:"error,omitempty"`
}

type FixUsersResponse struct {
	AppID       uint            `json:"app_id"`
	AppName     string          `json:"app_name"`
	TotalUsers  int             `json:"total_users"`
	EmailsFixed int             `json:"emails_fixed"`
	Subscribed  int             `json:"subscribed"`
	AlreadyOK   int             `json:"already_ok"`
	Failed      int             `json:"failed"`
	DryRun      bool            `json:"dry_run"`
	Results     []UserFixResult `json:"results"`
}

// FixUsers fixes users with malformed emails and creates subscriptions for unsubscribed users
func FixUsers(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	if appIDStr == "" {
		c.BadRequest("Missing app_id parameter")
		return
	}

	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		c.BadRequest("Invalid app_id parameter")
		return
	}

	var body FixUsersRequest
	_ = c.ShouldBindJSON(&body)

	// Verify the app belongs to dogfood account
	var app models.App
	if err := c.DB().FindForID(uint(appID), &app); err != nil {
		c.NotFound("App not found")
		return
	}

	var account models.Account
	if err := c.DB().FindForID(app.AccountID, &account); err != nil || account.Email != DogfoodAccountEmail {
		c.Unauthorized("Access denied - not a dogfood app")
		return
	}

	// Find the free plan for this app
	var freePlan models.Plan
	if err := c.DB().FindFirst(&freePlan, "app_id = ? AND is_free = ?", app.ID, true); err != nil {
		c.ServerError("Free plan not found for this app", err)
		return
	}

	// Get all users for this app
	var users []models.User
	if err := c.DB().List(&users, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list users", err)
		return
	}

	// Get existing subscriptions
	subMap := make(map[uint]bool)
	var subscriptions []models.Subscription
	if len(users) > 0 {
		userIDs := make([]uint, len(users))
		for i, u := range users {
			userIDs[i] = u.ID
		}
		c.DB().List(&subscriptions, "app_id = ? AND user_id IN ?", uint(appID), userIDs)
		for _, s := range subscriptions {
			subMap[s.UserID] = true
		}
	}

	response := FixUsersResponse{
		AppID:      uint(appID),
		AppName:    app.Name,
		TotalUsers: len(users),
		DryRun:     body.DryRun,
		Results:    make([]UserFixResult, 0),
	}

	for _, user := range users {
		result := UserFixResult{
			UserID: user.PublicID,
			Email:  user.Email,
			Name:   user.Name,
		}

		needsEmailFix := false
		needsSubscription := !subMap[user.ID]

		// Check if email needs fixing (no @ symbol means it's malformed)
		if !strings.Contains(user.Email, "@") {
			needsEmailFix = true
			result.OldEmail = user.Email
			result.NewEmail = fmt.Sprintf("%s@app.portcall.internal", user.Email)
		}

		// Determine action
		if needsEmailFix && needsSubscription {
			result.Action = "fix_email_and_subscribe"
		} else if needsEmailFix {
			result.Action = "fix_email"
		} else if needsSubscription {
			result.Action = "subscribe"
		} else {
			result.Action = "ok"
			response.AlreadyOK++
			continue // Don't add to results unless there's something to report
		}

		if body.DryRun {
			response.Results = append(response.Results, result)
			if needsEmailFix {
				response.EmailsFixed++
			}
			if needsSubscription {
				response.Subscribed++
			}
			continue
		}

		// Perform actual fixes
		hasError := false

		// Fix email if needed
		if needsEmailFix {
			newEmail := fmt.Sprintf("%s@app.portcall.internal", user.Email)
			if err := c.DB().Update(&user, "email", newEmail); err != nil {
				result.Error = fmt.Sprintf("failed to fix email: %v", err)
				hasError = true
			} else {
				log.Printf("[dogfood/fix] Fixed email for user %s: %s -> %s", user.PublicID, result.OldEmail, newEmail)
				result.Email = newEmail
				response.EmailsFixed++
			}
		}

		// Create subscription if needed
		if needsSubscription && !hasError {
			// Calculate next reset at end of current month
			now := time.Now().UTC()
			nextReset := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)

			sub := models.Subscription{
				PublicID:    dbx.GenPublicID("sub"),
				AppID:       app.ID,
				UserID:      user.ID,
				PlanID:      &freePlan.ID,
				Status:      "active",
				NextResetAt: nextReset,
			}
			if err := c.DB().Create(&sub); err != nil {
				result.Error = fmt.Sprintf("failed to create subscription: %v", err)
				hasError = true
			} else {
				log.Printf("[dogfood/fix] Created subscription for user %s to plan %s", user.PublicID, freePlan.Name)
				result.Subscribed = true
				result.PlanName = freePlan.Name
				response.Subscribed++
			}
		}

		if hasError {
			response.Failed++
		}

		response.Results = append(response.Results, result)
	}

	c.OK(response)
}
