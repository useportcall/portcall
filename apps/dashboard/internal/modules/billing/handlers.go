package billing

import (
	"fmt"
	"log"
	"os"

	"github.com/useportcall/portcall/libs/go/routerx"
	portcall "github.com/useportcall/portcall/sdks/go"
	"github.com/useportcall/portcall/sdks/go/resources"
)

var DFFeatures = struct {
	MaxSubscriptions string
	NumberOfUsers    string
}{
	MaxSubscriptions: "max_subscriptions",
	NumberOfUsers:    "number_of_users",
}

// QuotaStatus represents the usage and quota for a single feature.
// This matches the frontend's expected format.
type QuotaStatus struct {
	FeatureID   string `json:"feature_id"`
	Usage       int64  `json:"usage"`
	Quota       int64  `json:"quota"`
	Remaining   int64  `json:"remaining"`
	IsExceeded  bool   `json:"is_exceeded"`
	IsUnlimited bool   `json:"is_unlimited"`
}

// AllQuotasResponse contains quota information for all metered features.
// This matches the frontend's expected format.
type AllQuotasResponse struct {
	Subscriptions *QuotaStatus `json:"subscriptions,omitempty"`
	Users         *QuotaStatus `json:"users,omitempty"`
}

// CheckSubscriptionQuota is an HTTP handler to check subscription quota
func CheckSubscriptionQuota(ctx *routerx.Context) {
	client, err := getPortcall(ctx)
	if err != nil {
		log.Printf("[billing] Failed to get portcall client: %v", err)
		ctx.ServerError("Failed to get portcall client", err)
		return
	}

	entitlement, err := client.Entitlements.Get(ctx, ctx.PublicAppID(), DFFeatures.MaxSubscriptions)
	if err != nil {
		log.Printf("[billing] Failed to check subscription quota: %v", err)
		ctx.ServerError("Failed to check quota", err)
		return
	}

	data := new(QuotaStatus)
	data.FeatureID = DFFeatures.MaxSubscriptions
	data.Usage = entitlement.Usage
	data.Quota = entitlement.Quota

	if data.Quota == -1 {
		data.IsUnlimited = true
		data.IsExceeded = false
		data.Remaining = -1
	} else {
		data.Remaining = data.Quota - data.Usage
		if data.Remaining < 0 {
			data.IsExceeded = true
		}
	}

	ctx.OK(data)
}

// CheckUserQuota is an HTTP handler to check user quota
func CheckUserQuota(ctx *routerx.Context) {
	client, err := getPortcall(ctx)
	if err != nil {
		log.Printf("[billing] Failed to get portcall client: %v", err)
		ctx.ServerError("Failed to get portcall client", err)
		return
	}

	entitlement, err := client.Entitlements.Get(ctx, ctx.PublicAppID(), DFFeatures.NumberOfUsers)
	if err != nil {
		log.Printf("[billing] Failed to check user quota: %v", err)
		ctx.ServerError("Failed to check quota", err)
		return
	}

	data := new(QuotaStatus)
	data.FeatureID = DFFeatures.NumberOfUsers
	data.Usage = entitlement.Usage
	data.Quota = entitlement.Quota

	if data.Quota == -1 {
		data.IsUnlimited = true
		data.IsExceeded = false
		data.Remaining = -1
	} else {
		data.Remaining = data.Quota - data.Usage
		if data.Remaining < 0 {
			data.IsExceeded = true
		}
	}

	ctx.OK(data)
}

// UserBillingSubscriptionResponse contains the current app's plan and quota info.
// This matches the frontend's expected format.
type UserBillingSubscriptionResponse struct {
	PlanName         string                   `json:"plan_name"`
	IsFree           bool                     `json:"is_free"`
	HasPaymentMethod bool                     `json:"has_payment_method"`
	NextResetAt      *string                  `json:"next_reset_at,omitempty"`
	ScheduledPlanID  *string                  `json:"scheduled_plan_id,omitempty"`
	ScheduledPlan    *UserBillingPlanResponse `json:"scheduled_plan,omitempty"`
	CurrentPlan      *UserBillingPlanResponse `json:"current_plan,omitempty"`
}

// UserBillingPlanResponse represents a plan in billing responses
type UserBillingPlanResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getPortcall(ctx *routerx.Context) (*portcall.Client, error) {
	isLive := ctx.GetBool("is_live")

	var secret string
	if isLive {
		secret = os.Getenv("DOGFOOD_LIVE_SECRET")
	} else {
		secret = os.Getenv("DOGFOOD_TEST_SECRET")
	}

	if secret == "" {
		log.Println("Dogfood secret is not set, isLive:", isLive)
		return nil, fmt.Errorf("df secret is not set")
	}

	baseURL := os.Getenv("DF_API_URL")
	if baseURL == "" {
		baseURL = "https://api.useportcall.com"
	}

	return portcall.New(portcall.Config{BaseURL: baseURL, APIKey: secret}), nil
}

func getFreePlan(ctx *routerx.Context) (string, error) {
	isLive := ctx.GetBool("is_live")

	var plan string
	if isLive {
		plan = os.Getenv("DF_LIVE_FREE")
	} else {
		plan = os.Getenv("DF_TEST_FREE")
	}

	if plan == "" {
		log.Println("Dogfood free plan is not set, isLive:", isLive)
		return "", fmt.Errorf("df free plan is not set")
	}

	return plan, nil
}

func getProPlan(ctx *routerx.Context) (string, error) {
	isLive := ctx.GetBool("is_live")

	var plan string
	if isLive {
		plan = os.Getenv("DF_LIVE_PRO")
	} else {
		plan = os.Getenv("DF_TEST_PRO")
	}

	if plan == "" {
		log.Println("Dogfood pro plan is not set, isLive:", isLive)
		return "", fmt.Errorf("df pro plan is not set")
	}

	return plan, nil
}

// GetSubscriptionInfo is an HTTP handler to get subscription info for an app
func GetSubscriptionInfo(ctx *routerx.Context) {
	client, err := getPortcall(ctx)
	if err != nil {
		log.Printf("[billing] Failed to get portcall client: %v", err)
		ctx.ServerError("Failed to get portcall client", err)
		return
	}

	billingUserID := ctx.PublicAppID()
	subscription, err := client.Users.GetSubscription(ctx, billingUserID)
	if err != nil {
		log.Printf("[billing] Failed to get subscription info: %v", err)
		ctx.ServerError("Failed to get subscription info", err)
		return
	}

	planName := "Free"
	var currentPlan *UserBillingPlanResponse
	if subscription.Plan != nil {
		planName = subscription.Plan.Name
		currentPlan = &UserBillingPlanResponse{
			ID:   subscription.Plan.ID,
			Name: subscription.Plan.Name,
		}
	}

	freePlan, err := getFreePlan(ctx)
	if err != nil {
		log.Printf("[billing] Failed to get free plan: %v", err)
		ctx.ServerError("Failed to get free plan", err)
		return
	}

	isFree := subscription.Plan.ID == freePlan
	log.Printf("plan comparison: %s == %s", subscription.Plan.ID, freePlan)

	response := UserBillingSubscriptionResponse{
		PlanName:         planName,
		IsFree:           isFree,
		HasPaymentMethod: subscription.HasPaymentMethod,
		CurrentPlan:      currentPlan,
	}

	// Add next reset date if available
	if subscription.NextResetAt != nil {
		nextReset := subscription.NextResetAt.Format("2006-01-02T15:04:05Z07:00")
		response.NextResetAt = &nextReset
	}

	// Add scheduled plan if there's a pending change
	if subscription.ScheduledPlan != nil {
		response.ScheduledPlanID = &subscription.ScheduledPlan.ID
		response.ScheduledPlan = &UserBillingPlanResponse{
			ID:   subscription.ScheduledPlan.ID,
			Name: subscription.ScheduledPlan.Name,
		}
	}

	ctx.OK(response)
}

// UpgradeToProRequest is the request to upgrade to Pro
type UpgradeToProRequest struct {
	CancelURL   string `json:"cancel_url"`
	RedirectURL string `json:"redirect_url"`
}

// UpgradeToProResponse is the response with checkout URL or success status
type UpgradeToProResponse struct {
	CheckoutURL string `json:"checkout_url,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
	Success     bool   `json:"success"`
}

// UpgradeToPro is an HTTP handler to upgrade to Pro.
// If the user has a payment method, it updates the subscription directly.
// If not, it creates a checkout session for the user to add payment details.
func UpgradeToPro(c *routerx.Context) {
	var req UpgradeToProRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("[billing] Invalid upgrade request body: %v", err)
		c.BadRequest("Invalid request body")
		return
	}

	// Default URLs if not provided
	dashboardURL := os.Getenv("DASHBOARD_URL")
	if dashboardURL == "" {
		dashboardURL = "https://dashboard.useportcall.com"
	}

	cancelURL := req.CancelURL
	if cancelURL == "" {
		cancelURL = dashboardURL + "/usage?upgrade=cancelled"
	}

	redirectURL := req.RedirectURL
	if redirectURL == "" {
		redirectURL = dashboardURL + "/usage?upgrade=success"
	}

	client, err := getPortcall(c)
	if err != nil {
		log.Printf("[billing] Failed to get portcall client: %v", err)
		c.ServerError("Failed to get portcall client", err)
		return
	}

	billingUserID := c.PublicAppID()

	// Get the user's current subscription to check for payment method
	subscription, err := client.Users.GetSubscription(c, billingUserID)
	if err != nil {
		log.Printf("[billing] Failed to get subscription info: %v", err)
		c.ServerError("Failed to get subscription info", err)
		return
	}

	planID, err := getProPlan(c)
	if err != nil {
		log.Printf("[billing] Failed to get pro plan: %v", err)
		c.ServerError("Failed to get pro plan", err)
		return
	}

	// If user has a payment method, update subscription directly
	if subscription.HasPaymentMethod {
		log.Printf("[billing] User %s has payment method, updating subscription directly", billingUserID)

		_, err = client.Subscriptions.Update(c, subscription.ID, portcall.UpdateSubscriptionInput{PlanID: &planID})
		if err != nil {
			log.Printf("[billing] Failed to update subscription: %v", err)
			c.ServerError("Failed to update subscription", err)
			return
		}

		log.Printf("[billing] Successfully upgraded subscription %s to Pro tier", subscription.ID)
		c.OK(UpgradeToProResponse{Success: true})
		return
	}

	// No payment method, create checkout session
	log.Printf("[billing] Creating checkout session for app %s", billingUserID)

	session, err := client.CheckoutSessions.Create(c, resources.CreateCheckoutSessionRequest{
		PlanID:      planID,
		UserID:      c.PublicAppID(),
		CancelURL:   cancelURL,
		RedirectURL: redirectURL,
	})
	if err != nil {
		log.Printf("[billing] Failed to create checkout session: %v", err)
		c.ServerError("Failed to create checkout session", err)
		return
	}

	c.OK(UpgradeToProResponse{
		CheckoutURL: session.URL,
		SessionID:   session.ID,
		Success:     false,
	})
}

// DowngradeToFreeRequest is the request to downgrade to Free
type DowngradeToFreeRequest struct {
	Immediate bool `json:"immediate"` // If true, apply immediately. Default: schedule for next reset.
}

// DowngradeToFreeResponse is the response for downgrade
type DowngradeToFreeResponse struct {
	Success   bool   `json:"success"`
	Scheduled bool   `json:"scheduled"` // If true, downgrade is scheduled for next reset
	Message   string `json:"message,omitempty"`
}

// DowngradeToFree is an HTTP handler to downgrade subscription to Free tier.
// By default, the downgrade is scheduled for the next reset date.
// Pass immediate=true to apply immediately.
func DowngradeToFree(c *routerx.Context) {
	var req DowngradeToFreeRequest
	if err := c.BindJSON(&req); err != nil {
		// Non-JSON request is fine, just use defaults
		req = DowngradeToFreeRequest{Immediate: false}
	}

	client, err := getPortcall(c)
	if err != nil {
		log.Printf("[billing] Failed to get portcall client: %v", err)
		c.ServerError("Failed to get portcall client", err)
		return
	}

	billingUserID := c.PublicAppID()

	// Get the user's current subscription
	subscription, err := client.Users.GetSubscription(c, billingUserID)
	if err != nil {
		log.Printf("[billing] Failed to get subscription info: %v", err)
		c.ServerError("Failed to get subscription info", err)
		return
	}

	freePlan, err := getFreePlan(c)
	if err != nil {
		log.Printf("[billing] Failed to get free plan: %v", err)
		c.ServerError("Failed to get free plan", err)
		return
	}

	// Check if already on free tier
	if subscription.Plan != nil && subscription.Plan.ID == freePlan {
		log.Printf("[billing] User %s is already on free tier", billingUserID)
		c.BadRequest("Already on free tier")
		return
	}

	applyAtNextReset := !req.Immediate

	log.Printf("[billing] Downgrading subscription %s to Free tier for user %s (scheduled=%v)", subscription.ID, billingUserID, applyAtNextReset)

	_, err = client.Subscriptions.Update(c, subscription.ID, portcall.UpdateSubscriptionInput{
		PlanID:           &freePlan,
		ApplyAtNextReset: &applyAtNextReset,
	})
	if err != nil {
		log.Printf("[billing] Failed to update subscription: %v", err)
		c.ServerError("Failed to downgrade subscription", err)
		return
	}

	message := "Successfully downgraded to Free tier"
	if applyAtNextReset {
		message = "Downgrade to Free tier scheduled for next billing cycle"
	}

	log.Printf("[billing] Successfully downgraded subscription %s to Free tier (scheduled=%v)", subscription.ID, applyAtNextReset)
	c.OK(DowngradeToFreeResponse{
		Success:   true,
		Scheduled: applyAtNextReset,
		Message:   message,
	})
}

// BillingInvoice represents an invoice for the billing view
type BillingInvoice struct {
	ID            string  `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	Currency      string  `json:"currency"`
	Total         int64   `json:"total"`
	Status        string  `json:"status"`
	PDFURL        *string `json:"pdf_url,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// ListBillingInvoicesResponse is the response for listing billing invoices
type ListBillingInvoicesResponse struct {
	Invoices []BillingInvoice `json:"invoices"`
}

// ListBillingInvoices is an HTTP handler to list invoices for the dashboard app billing
func ListBillingInvoices(c *routerx.Context) {
	client, err := getPortcall(c)
	if err != nil {
		log.Printf("[billing] Failed to get portcall client: %v", err)
		c.ServerError("Failed to get portcall client", err)
		return
	}

	billingUserID := c.PublicAppID()

	invoices, err := client.Invoices.List(c, &portcall.ListInvoicesParams{UserID: billingUserID})
	if err != nil {
		log.Printf("[billing] Failed to list invoices: %v", err)
		c.ServerError("Failed to list invoices", err)
		return
	}

	response := ListBillingInvoicesResponse{
		Invoices: make([]BillingInvoice, len(invoices)),
	}

	for i, inv := range invoices {
		response.Invoices[i] = BillingInvoice{
			ID:            inv.ID,
			InvoiceNumber: inv.InvoiceNumber,
			Currency:      inv.Currency,
			Total:         inv.Total,
			Status:        inv.Status,
			PDFURL:        inv.PDFURL,
			CreatedAt:     inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.OK(response)
}

// BillingAddressResponse represents the billing address response
type BillingAddressResponse struct {
	ID         string `json:"id"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// GetBillingAddressResponse wraps the billing address response
type GetBillingAddressResponse struct {
	BillingAddress *BillingAddressResponse `json:"billing_address"`
}

// GetBillingAddress is an HTTP handler to get the billing address for the dashboard app
func GetBillingAddress(c *routerx.Context) {
	client, err := getPortcall(c)
	if err != nil {
		log.Printf("[billing] Failed to get portcall client: %v", err)
		c.ServerError("Failed to get portcall client", err)
		return
	}

	billingUserID := c.PublicAppID()

	user, err := client.Users.Get(c, billingUserID)
	if err != nil {
		log.Printf("[billing] Failed to get user: %v", err)
		c.ServerError("Failed to get user", err)
		return
	}

	response := GetBillingAddressResponse{}

	if user.BillingAddress != nil {
		response.BillingAddress = &BillingAddressResponse{
			ID:         user.BillingAddress.ID,
			Line1:      user.BillingAddress.Line1,
			Line2:      user.BillingAddress.Line2,
			City:       user.BillingAddress.City,
			State:      user.BillingAddress.State,
			PostalCode: user.BillingAddress.PostalCode,
			Country:    user.BillingAddress.Country,
		}
	}

	c.OK(response)
}

// UpsertBillingAddressRequest is the request to create or update a billing address
type UpsertBillingAddressRequest struct {
	Line1      string `json:"line1" binding:"required"`
	Line2      string `json:"line2"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code" binding:"required"`
	Country    string `json:"country" binding:"required"`
}

// UpsertBillingAddressHandler is an HTTP handler to create or update the billing address
func UpsertBillingAddressHandler(c *routerx.Context) {
	var req UpsertBillingAddressRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("[billing] Invalid billing address request body: %v", err)
		c.BadRequest("Invalid request body")
		return
	}

	billingUserID := c.PublicAppID()

	address, err := UpsertBillingAddress(c, billingUserID, portcall.UpsertBillingAddressInput{
		Line1:      req.Line1,
		Line2:      req.Line2,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
	})
	if err != nil {
		log.Printf("[billing] Failed to upsert billing address: %v", err)
		c.ServerError("Failed to save billing address", err)
		return
	}

	response := BillingAddressResponse{
		ID:         address.ID,
		Line1:      address.Line1,
		Line2:      address.Line2,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
	}

	c.OK(response)
}
