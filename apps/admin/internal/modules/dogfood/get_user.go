package dogfood

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// UserDetail represents detailed information about a dogfood user
type UserDetail struct {
	ID           uint                `json:"id"`
	PublicID     string              `json:"public_id"`
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	CreatedAt    string              `json:"created_at"`
	Subscription *SubscriptionDetail `json:"subscription,omitempty"`
	Entitlements []EntitlementInfo   `json:"entitlements"`
	Invoices     []InvoiceInfo       `json:"invoices"`
}

type SubscriptionDetail struct {
	ID                   uint   `json:"id"`
	PublicID             string `json:"public_id"`
	Status               string `json:"status"`
	PlanID               *uint  `json:"plan_id,omitempty"`
	PlanName             string `json:"plan_name,omitempty"`
	Currency             string `json:"currency"`
	BillingInterval      string `json:"billing_interval"`
	BillingIntervalCount int    `json:"billing_interval_count"`
	LastResetAt          string `json:"last_reset_at"`
	NextResetAt          string `json:"next_reset_at"`
	CreatedAt            string `json:"created_at"`
}

type EntitlementInfo struct {
	ID              uint   `json:"id"`
	FeaturePublicID string `json:"feature_public_id"`
	Usage           int64  `json:"usage"`
	Quota           int64  `json:"quota"`
	Interval        string `json:"interval"`
	IsMetered       bool   `json:"is_metered"`
	LastResetAt     string `json:"last_reset_at"`
	NextResetAt     string `json:"next_reset_at"`
}

type InvoiceInfo struct {
	ID        uint   `json:"id"`
	PublicID  string `json:"public_id"`
	Status    string `json:"status"`
	Total     int64  `json:"total"`
	Currency  string `json:"currency"`
	CreatedAt string `json:"created_at"`
	DueBy     string `json:"due_by"`
}

// GetUser gets detailed information about a dogfood user
func GetUser(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	userID := c.Param("user_id")

	if appIDStr == "" {
		c.BadRequest("Missing app_id parameter")
		return
	}
	if userID == "" {
		c.BadRequest("Missing user_id parameter")
		return
	}

	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		c.BadRequest("Invalid app_id parameter")
		return
	}

	// Verify the app belongs to dogfood account
	var app models.App
	if err := c.DB().FindForID(uint(appID), &app); err != nil {
		c.NotFound("App not found")
		return
	}

	var account models.Account
	if err := c.DB().FindForID(app.AccountID, &account); err != nil || account.Email != DogfoodAccountEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied - not a dogfood app"})
		return
	}

	// Get user
	var user models.User
	if err := c.DB().FindFirst(&user, "app_id = ? AND public_id = ?", uint(appID), userID); err != nil {
		c.NotFound("User not found")
		return
	}

	result := UserDetail{
		ID:           user.ID,
		PublicID:     user.PublicID,
		Name:         user.Name,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Entitlements: []EntitlementInfo{},
		Invoices:     []InvoiceInfo{},
	}

	// Get subscription
	var subscription models.Subscription
	if err := c.DB().FindFirst(&subscription, "app_id = ? AND user_id = ?", uint(appID), user.ID); err == nil {
		planName := ""
		if subscription.PlanID != nil {
			var plan models.Plan
			if err := c.DB().FindForID(*subscription.PlanID, &plan); err == nil {
				planName = plan.Name
			}
		}

		result.Subscription = &SubscriptionDetail{
			ID:                   subscription.ID,
			PublicID:             subscription.PublicID,
			Status:               subscription.Status,
			PlanID:               subscription.PlanID,
			PlanName:             planName,
			Currency:             subscription.Currency,
			BillingInterval:      subscription.BillingInterval,
			BillingIntervalCount: subscription.BillingIntervalCount,
			LastResetAt:          subscription.LastResetAt.Format("2006-01-02T15:04:05Z"),
			NextResetAt:          subscription.NextResetAt.Format("2006-01-02T15:04:05Z"),
			CreatedAt:            subscription.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	// Get entitlements
	var entitlements []models.Entitlement
	if err := c.DB().List(&entitlements, "app_id = ? AND user_id = ?", uint(appID), user.ID); err == nil {
		for _, e := range entitlements {
			lastReset := ""
			nextReset := ""
			if e.LastResetAt != nil {
				lastReset = e.LastResetAt.Format("2006-01-02T15:04:05Z")
			}
			if e.NextResetAt != nil {
				nextReset = e.NextResetAt.Format("2006-01-02T15:04:05Z")
			}

			result.Entitlements = append(result.Entitlements, EntitlementInfo{
				ID:              e.ID,
				FeaturePublicID: e.FeaturePublicID,
				Usage:           e.Usage,
				Quota:           e.Quota,
				Interval:        e.Interval,
				IsMetered:       e.IsMetered,
				LastResetAt:     lastReset,
				NextResetAt:     nextReset,
			})
		}
	}

	// Get invoices
	var invoices []models.Invoice
	if err := c.DB().List(&invoices, "app_id = ? AND user_id = ? ORDER BY created_at DESC", uint(appID), user.ID); err == nil {
		for _, i := range invoices {
			result.Invoices = append(result.Invoices, InvoiceInfo{
				ID:        i.ID,
				PublicID:  i.PublicID,
				Status:    i.Status,
				Total:     i.Total,
				Currency:  i.Currency,
				CreatedAt: i.CreatedAt.Format("2006-01-02T15:04:05Z"),
				DueBy:     i.DueBy.Format("2006-01-02T15:04:05Z"),
			})
		}
	}

	c.OK(result)
}
