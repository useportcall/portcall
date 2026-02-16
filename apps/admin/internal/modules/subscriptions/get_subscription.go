package subscriptions

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type SubscriptionDetail struct {
	ID              uint            `json:"id"`
	PublicID        string          `json:"public_id"`
	UserID          uint            `json:"user_id"`
	UserName        string          `json:"user_name"`
	UserEmail       string          `json:"user_email"`
	PlanID          *uint           `json:"plan_id"`
	PlanName        string          `json:"plan_name"`
	Status          string          `json:"status"`
	BillingInterval string          `json:"billing_interval"`
	Currency        string          `json:"currency"`
	NextResetAt     time.Time       `json:"next_reset_at"`
	LastResetAt     time.Time       `json:"last_reset_at"`
	Items           []ItemDetail    `json:"items"`
	RecentInvoices  []InvoiceDetail `json:"recent_invoices"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ItemDetail struct {
	ID           uint   `json:"id"`
	PublicID     string `json:"public_id"`
	Title        string `json:"title"`
	Quantity     int32  `json:"quantity"`
	UnitAmount   int64  `json:"unit_amount"`
	PricingModel string `json:"pricing_model"`
}

type InvoiceDetail struct {
	ID        uint      `json:"id"`
	PublicID  string    `json:"public_id"`
	Total     int64     `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func GetSubscription(c *routerx.Context) {
	subIDStr := c.Param("subscription_id")
	subID, err := strconv.ParseUint(subIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid subscription ID")
		return
	}

	var sub models.Subscription
	if err := c.DB().FindForID(uint(subID), &sub); err != nil {
		c.NotFound("Subscription not found")
		return
	}

	// Get user info
	var user models.User
	userName := ""
	userEmail := ""
	if err := c.DB().FindForID(sub.UserID, &user); err == nil {
		userName = user.Name
		userEmail = user.Email
	}

	// Get plan name
	planName := ""
	if sub.PlanID != nil {
		var plan models.Plan
		if err := c.DB().FindForID(*sub.PlanID, &plan); err == nil {
			planName = plan.Name
		}
	}

	// Get items
	var items []models.SubscriptionItem
	c.DB().List(&items, "subscription_id = ?", sub.ID)

	itemDetails := make([]ItemDetail, len(items))
	for i, item := range items {
		itemDetails[i] = ItemDetail{
			ID:           item.ID,
			PublicID:     item.PublicID,
			Title:        item.Title,
			Quantity:     item.Quantity,
			UnitAmount:   item.UnitAmount,
			PricingModel: item.PricingModel,
		}
	}

	// Get invoices
	var invoices []models.Invoice
	c.DB().ListWithOrderAndLimit(&invoices, "created_at DESC", 10, "subscription_id = ?", sub.ID)

	invoiceDetails := make([]InvoiceDetail, len(invoices))
	for i, inv := range invoices {
		invoiceDetails[i] = InvoiceDetail{
			ID:        inv.ID,
			PublicID:  inv.PublicID,
			Total:     inv.Total,
			Status:    inv.Status,
			CreatedAt: inv.CreatedAt,
		}
	}

	result := SubscriptionDetail{
		ID:              sub.ID,
		PublicID:        sub.PublicID,
		UserID:          sub.UserID,
		UserName:        userName,
		UserEmail:       userEmail,
		PlanID:          sub.PlanID,
		PlanName:        planName,
		Status:          sub.Status,
		BillingInterval: sub.BillingInterval,
		Currency:        sub.Currency,
		NextResetAt:     sub.NextResetAt,
		LastResetAt:     sub.LastResetAt,
		Items:           itemDetails,
		RecentInvoices:  invoiceDetails,
		CreatedAt:       sub.CreatedAt,
		UpdatedAt:       sub.UpdatedAt,
	}

	c.OK(result)
}
