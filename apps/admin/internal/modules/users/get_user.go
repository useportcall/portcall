package users

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UserDetail struct {
	ID               uint                  `json:"id"`
	PublicID         string                `json:"public_id"`
	Name             string                `json:"name"`
	Email            string                `json:"email"`
	HasSubscription  bool                  `json:"has_subscription"`
	HasPaymentMethod bool                  `json:"has_payment_method"`
	Subscriptions    []SubscriptionSummary `json:"subscriptions"`
	Entitlements     []EntitlementSummary  `json:"entitlements"`
	Invoices         []InvoiceSummary      `json:"invoices"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

type SubscriptionSummary struct {
	ID       uint   `json:"id"`
	PublicID string `json:"public_id"`
	PlanName string `json:"plan_name"`
	Status   string `json:"status"`
}

type EntitlementSummary struct {
	ID              uint   `json:"id"`
	FeaturePublicID string `json:"feature_public_id"`
	Quota           int64  `json:"quota"`
	Usage           int64  `json:"usage"`
}

type InvoiceSummary struct {
	ID       uint   `json:"id"`
	PublicID string `json:"public_id"`
	Total    int64  `json:"total"`
	Status   string `json:"status"`
}

func GetUser(c *routerx.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid user ID")
		return
	}

	var user models.User
	if err := c.DB().FindForID(uint(userID), &user); err != nil {
		c.NotFound("User not found")
		return
	}

	// Get subscriptions
	var subscriptions []models.Subscription
	c.DB().List(&subscriptions, "user_id = ?", user.ID)

	subSummaries := make([]SubscriptionSummary, len(subscriptions))
	for i, sub := range subscriptions {
		planName := ""
		if sub.PlanID != nil {
			var plan models.Plan
			if err := c.DB().FindForID(*sub.PlanID, &plan); err == nil {
				planName = plan.Name
			}
		}
		subSummaries[i] = SubscriptionSummary{
			ID:       sub.ID,
			PublicID: sub.PublicID,
			PlanName: planName,
			Status:   sub.Status,
		}
	}

	// Get entitlements
	var entitlements []models.Entitlement
	c.DB().List(&entitlements, "user_id = ?", user.ID)

	entSummaries := make([]EntitlementSummary, len(entitlements))
	for i, ent := range entitlements {
		entSummaries[i] = EntitlementSummary{
			ID:              ent.ID,
			FeaturePublicID: ent.FeaturePublicID,
			Quota:           ent.Quota,
			Usage:           ent.Usage,
		}
	}

	// Get invoices
	var invoices []models.Invoice
	c.DB().ListWithOrderAndLimit(&invoices, "created_at DESC", 10, "user_id = ?", user.ID)

	invSummaries := make([]InvoiceSummary, len(invoices))
	for i, inv := range invoices {
		invSummaries[i] = InvoiceSummary{
			ID:       inv.ID,
			PublicID: inv.PublicID,
			Total:    inv.Total,
			Status:   inv.Status,
		}
	}

	var pmCount int64
	c.DB().Count(&pmCount, models.PaymentMethod{}, "user_id = ?", user.ID)

	result := UserDetail{
		ID:               user.ID,
		PublicID:         user.PublicID,
		Name:             user.Name,
		Email:            user.Email,
		HasSubscription:  len(subscriptions) > 0,
		HasPaymentMethod: pmCount > 0,
		Subscriptions:    subSummaries,
		Entitlements:     entSummaries,
		Invoices:         invSummaries,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}

	c.OK(result)
}
