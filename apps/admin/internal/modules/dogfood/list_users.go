package dogfood

import (
	"strconv"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// UserWithSubscription represents a dogfood user with their subscription info
type UserWithSubscription struct {
	ID           uint               `json:"id"`
	PublicID     string             `json:"public_id"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	CreatedAt    string             `json:"created_at"`
	Subscription *SubscriptionBrief `json:"subscription,omitempty"`
}

type SubscriptionBrief struct {
	ID          uint   `json:"id"`
	PublicID    string `json:"public_id"`
	Status      string `json:"status"`
	PlanName    string `json:"plan_name,omitempty"`
	NextResetAt string `json:"next_reset_at"`
}

// ListUsers lists all users for a dogfood app
func ListUsers(c *routerx.Context) {
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

	// Get all users for this app
	var users []models.User
	if err := c.DB().List(&users, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list users", err)
		return
	}

	// Get subscriptions for these users
	var subscriptions []models.Subscription
	if len(users) > 0 {
		userIDs := make([]uint, len(users))
		for i, u := range users {
			userIDs[i] = u.ID
		}
		c.DB().List(&subscriptions, "app_id = ? AND user_id IN ?", uint(appID), userIDs)
	}

	// Get plans for subscriptions
	planMap := make(map[uint]string)
	if len(subscriptions) > 0 {
		planIDs := make([]uint, 0)
		for _, s := range subscriptions {
			if s.PlanID != nil {
				planIDs = append(planIDs, *s.PlanID)
			}
		}
		if len(planIDs) > 0 {
			var plans []models.Plan
			c.DB().List(&plans, "id IN ?", planIDs)
			for _, p := range plans {
				planMap[p.ID] = p.Name
			}
		}
	}

	// Build subscription map
	subMap := make(map[uint]*SubscriptionBrief)
	for _, s := range subscriptions {
		planName := ""
		if s.PlanID != nil {
			planName = planMap[*s.PlanID]
		}
		subMap[s.UserID] = &SubscriptionBrief{
			ID:          s.ID,
			PublicID:    s.PublicID,
			Status:      s.Status,
			PlanName:    planName,
			NextResetAt: s.NextResetAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	// Build response
	result := make([]UserWithSubscription, len(users))
	for i, u := range users {
		result[i] = UserWithSubscription{
			ID:           u.ID,
			PublicID:     u.PublicID,
			Name:         u.Name,
			Email:        u.Email,
			CreatedAt:    u.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Subscription: subMap[u.ID],
		}
	}

	c.OK(result)
}
