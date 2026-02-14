package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/qx"
)

type IMultiplexer interface {
	HandleFunc(taskType string, handler HandlerFunc)
}

type multiplexer struct {
	instance *asynq.ServeMux
	db       dbx.IORM
	queue    qx.IQueue
	crypto   cryptox.ICrypto
	email    emailx.IEmailClient
}

func (m *multiplexer) HandleFunc(taskType string, handler HandlerFunc) {
	m.instance.HandleFunc(taskType, func(ctx context.Context, t *asynq.Task) error {
		c := &Context{
			Task:    t,
			orm:     m.db,     // Use the DB connection from the multiplexer
			queue:   m.queue,  // Use the Queue from the multiplexer
			crypto:  m.crypto, // Use the Crypto from the multiplexer
			email:   m.email,  // Use the Email client from the multiplexer
			Context: ctx,      // Pass the asynq context to the embedded Context field
		}

		log.Println("Processing task:", taskType)

		return handler(c)
	})
}

type HandlerFunc = func(IContext) error

type IContext interface {
	DB() dbx.IORM
	Crypto() cryptox.ICrypto
	Queue() qx.IQueue
	Payload() []byte
	EmailClient() emailx.IEmailClient
	GetDFSecret(isTest bool) (string, error)
	CheckInvoiceIdempotency(subscriptionID uint) (bool, error)
	CountInvoices(appID uint) (int64, error)
	ListSubscriptionItemIDs(subscription_id uint) ([]uint, error)
	BuildInvoice(*models.Subscription, *models.Company, int64) *models.Invoice
	EnqueueResolveCheckoutSession(id string, paymentMethodID string) error
	GetResolveCheckoutSessionPayload() (*ResolveCheckoutSessionPayload, error)
	FindCheckoutSessionByExternalID(externalID string) (*models.CheckoutSession, error)
	SaveCheckoutSession(*models.CheckoutSession) error
	EnqueueCreatePaymentMethod(appID uint, userID uint, planID uint, externalPaymentMethodID string) error
	GetCreatePaymentMethodPayload() (*CreatePaymentMethodPayload, error)
	BuildPaymentMethod(*CreatePaymentMethodPayload) *models.PaymentMethod
	UpsertPaymentMethod(pm *models.PaymentMethod) error
	EnqueueUpsertSubscription(appID uint, userID uint, planID uint) error
	GetUpsertSubscriptionPayload() (*UpsertSubscriptionPayload, error)
	IUpsertEntitlementsContext
	ICreateSubscriptionContext
	context.Context
}

type ICreateSubscriptionContext interface {
	GetCreateSubscriptionPayload() (*CreateSubscriptionPayload, error)
	FindUserBillingAddressID(userID uint) (*uint, error)
	BuildSubscription(*models.Plan) *models.Subscription
	SetSubscriptionRollback(*models.Subscription, *models.Plan) error
	SaveSubscription(*models.Subscription) error
}

type CreateSubscriptionPayload struct {
	PlanID uint `json:"plan_id"`
	UserID uint `json:"user_id"`
}

func (c *Context) GetCreateSubscriptionPayload() (*CreateSubscriptionPayload, error) {
	var p CreateSubscriptionPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		log.Printf("Error unmarshaling CreateSubscriptionPayload: %v", err)
		return nil, err
	}

	log.Printf("CreateSubscriptionPayload: PlanID=%d, UserID=%d", p.PlanID, p.UserID)

	return &p, nil
}

func (c *Context) FindUserBillingAddressID(userID uint) (*uint, error) {
	var user models.User
	if err := c.DB().FindForID(userID, &user); err != nil {
		log.Printf("Error finding user with ID %d: %v", userID, err)
		return nil, err
	}

	log.Printf("User ID %d has BillingAddressID: %v", userID, user.BillingAddressID)

	return user.BillingAddressID, nil
}

func (c *Context) BuildSubscription(plan *models.Plan) *models.Subscription {
	subscription := &models.Subscription{
		PublicID:             dbx.GenPublicID("sub"),
		AppID:                plan.AppID,
		Status:               "active",
		Currency:             plan.Currency,
		PlanID:               &plan.ID, // TODO: rollback plan?
		BillingInterval:      plan.Interval,
		BillingIntervalCount: plan.IntervalCount,
		InvoiceDueByDays:     plan.InvoiceDueByDays,
	}

	return subscription
}

// TODO: improve logic
func (c *Context) SetSubscriptionRollback(dest *models.Subscription, src *models.Plan) error {
	if !src.IsFree {
		freePlans := []models.Plan{}
		if err := c.DB().List(&freePlans, "app_id = ? AND is_free = ?", src.AppID, true); err != nil {
			log.Printf("Error listing free plans for AppID %d: %v", src.AppID, err)
			return err
		}

		if len(freePlans) > 0 {
			dest.RollbackPlanID = &freePlans[0].ID
			log.Printf("Set rollback plan ID %d for subscription", freePlans[0].ID)
		} else {
			log.Printf("No free plans found for AppID %d to set as rollback", src.AppID)
		}
	}
	return nil
}

func (c *Context) SaveSubscription(subscription *models.Subscription) error {
	if err := c.DB().Save(subscription); err != nil {
		log.Printf("Error saving subscription ID %d: %v", subscription.ID, err)
		return err
	}

	log.Printf("Subscription ID %d saved successfully", subscription.ID)
	return nil
}

type Context struct {
	Task   *asynq.Task
	orm    dbx.IORM
	queue  qx.IQueue
	crypto cryptox.ICrypto
	email  emailx.IEmailClient
	context.Context
}

func (c *Context) DB() dbx.IORM {
	return c.orm
}

func (c *Context) Queue() qx.IQueue {
	return c.queue
}

func (c *Context) Payload() []byte {
	return c.Task.Payload()
}

func (c *Context) Crypto() cryptox.ICrypto {
	return c.crypto
}

func (c *Context) EmailClient() emailx.IEmailClient {
	return c.email
}

func (c *Context) GetDFSecret(isTest bool) (string, error) {
	var secret string
	if isTest {
		secret = os.Getenv("DOGFOOD_TEST_SECRET")
	} else {
		secret = os.Getenv("DOGFOOD_LIVE_SECRET")
	}

	if secret == "" {
		log.Println("Dogfood secret is not set, isTest:", isTest)
		return "", fmt.Errorf("Dogfood secret is not set")
	}

	return secret, nil
}

// Idempotency check: prevent duplicate invoices for the same subscription in the same billing period
// Check if there's already a pending or issued invoice for this subscription since the last reset
func (c *Context) CheckInvoiceIdempotency(subscriptionID uint) (bool, error) {
	var existingInvoice models.Invoice
	if err := c.DB().FindFirst(&existingInvoice, "subscription_id = ? AND status IN (?, ?, ?)",
		subscriptionID, "pending", "issued", "paid"); err == nil {
		log.Printf("Invoice already exists for subscription ID %d, skipping creation", subscriptionID)
		// Invoice already exists for this subscription, skip creation
		return true, nil
	} else if !dbx.IsRecordNotFoundError(err) {
		log.Printf("Error checking existing invoices for subscription ID %d: %v", subscriptionID, err)
		return false, err
	}

	return false, nil
}

func (c *Context) CountInvoices(appID uint) (int64, error) {
	var count int64
	if err := c.DB().Count(&count, &models.Invoice{}, "app_id = ?", appID); err != nil {
		log.Printf("Error counting invoices for app ID %d: %v", appID, err)
		return 0, err
	}
	return count, nil
}

func (c *Context) ListSubscriptionItemIDs(subscription_id uint) ([]uint, error) {
	ids := []uint{}
	if err := c.DB().ListIDs("subscription_items", &ids, "subscription_id = ?", subscription_id); err != nil {
		log.Printf("Error listing subscription item IDs for subscription ID %d: %v", subscription_id, err)
		return nil, err
	}
	return ids, nil
}

func (c *Context) BuildInvoice(
	subscription *models.Subscription,
	company *models.Company,
	count int64,
) *models.Invoice {
	invoiceAppURL := os.Getenv("INVOICE_APP_URL")
	if invoiceAppURL == "" {
		// return fmt.Errorf("INVOICE_APP_URL environment variable is not set")
	}

	discountPct := 0
	if subscription.DiscountPct > 0 && count+1 <= int64(subscription.DiscountQty) {
		discountPct = subscription.DiscountPct
	}

	publicID := dbx.GenPublicID("invoice")

	return &models.Invoice{
		AppID:              subscription.AppID,
		SubscriptionID:     &subscription.ID,
		UserID:             subscription.UserID,
		PublicID:           publicID,
		Status:             "pending",
		Currency:           subscription.Currency,
		PDFURL:             fmt.Sprintf("%s/invoices/%s/view", invoiceAppURL, publicID),
		EmailURL:           fmt.Sprintf("%s/invoice-email/%s", invoiceAppURL, publicID),
		DueBy:              time.Now().AddDate(0, 0, subscription.InvoiceDueByDays),
		InvoiceNumber:      fmt.Sprintf("INV-%07d", count+1), // invoice number should be INV-0000001 format
		InvoiceNumberCount: count + 1,
		CompanyAddressID:   company.BillingAddressID,
		BillingAddressID:   *subscription.BillingAddressID,
		ShippingAddressID:  subscription.BillingAddressID,
		DiscountPct:        discountPct,
	}
}

// resolve checkout session task
type ResolveCheckoutSessionPayload struct {
	ExternalSessionID       string `json:"external_session_id"`
	ExternalPaymentMethodID string `json:"external_payment_method_id"`
}

func (c *Context) GetResolveCheckoutSessionPayload() (*ResolveCheckoutSessionPayload, error) {
	var p ResolveCheckoutSessionPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Context) FindCheckoutSessionByExternalID(externalID string) (*models.CheckoutSession, error) {
	var session models.CheckoutSession
	if err := c.DB().FindFirst(&session, "external_id = ?", externalID); err != nil {
		log.Printf("Error finding checkout session with external ID %s: %v", externalID, err)
		return nil, err
	}

	log.Printf("Found Checkout Session: ID=%d, AppID=%d, UserID=%d, PlanID=%d, Status=%s",
		session.ID, session.AppID, session.UserID, session.PlanID, session.Status)

	return &session, nil
}

func (c *Context) SaveCheckoutSession(session *models.CheckoutSession) error {
	if err := c.DB().Save(session); err != nil {
		log.Printf("Error saving checkout session ID %d: %v", session.ID, err)
		return err
	}
	return nil
}

func (c *Context) EnqueueResolveCheckoutSession(id string, paymentMethodID string) error {
	payload := ResolveCheckoutSessionPayload{
		ExternalSessionID:       id,
		ExternalPaymentMethodID: paymentMethodID,
	}

	if err := c.Queue().Enqueue("resolve_checkout_session", payload, "billing_queue"); err != nil {
		log.Printf("Error enqueueing resolve_checkout_session task: %v", err)
		return err
	}

	log.Printf("Enqueued resolve_checkout_session task for session ID %s", id)

	return nil
}

// create payment method task

type CreatePaymentMethodPayload struct {
	AppID                   uint   `json:"app_id"`
	UserID                  uint   `json:"user_id"`
	PlanID                  uint   `json:"plan_id"`
	ExternalPaymentMethodID string `json:"external_payment_method_id"`
}

func (c *Context) GetCreatePaymentMethodPayload() (*CreatePaymentMethodPayload, error) {
	var p CreatePaymentMethodPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Context) EnqueueCreatePaymentMethod(appID uint, userID uint, planID uint, externalPaymentMethodID string) error {
	payload := CreatePaymentMethodPayload{
		AppID:                   appID,
		UserID:                  userID,
		PlanID:                  planID,
		ExternalPaymentMethodID: externalPaymentMethodID,
	}

	if err := c.Queue().Enqueue("create_payment_method", payload, "billing_queue"); err != nil {
		log.Printf("Error enqueueing create_payment_method task: %v", err)
		return err
	}

	log.Printf("Enqueued create_payment_method task for user ID %d", userID)

	return nil
}

func (c *Context) BuildPaymentMethod(payload *CreatePaymentMethodPayload) *models.PaymentMethod {
	paymentMethod := models.PaymentMethod{}
	paymentMethod.PublicID = dbx.GenPublicID("pm")
	paymentMethod.AppID = payload.AppID
	paymentMethod.UserID = payload.UserID
	paymentMethod.ExternalID = payload.ExternalPaymentMethodID
	paymentMethod.ExternalType = "card"
	return &paymentMethod
}

func (c *Context) UpsertPaymentMethod(pm *models.PaymentMethod) error {
	if err := c.DB().FindFirst(&pm, "user_id = ? AND external_id = ?", pm.UserID, pm.ExternalID); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			return err
		}

		if err := c.DB().Create(&pm); err != nil {
			return err
		}

		log.Printf("Payment method created with ID %d for user ID %d", pm.ID, pm.UserID)
	}
	return nil
}

// upsert subscription task

type UpsertSubscriptionPayload struct {
	AppID  uint `json:"app_id"`
	UserID uint `json:"user_id"`
	PlanID uint `json:"plan_id"`
}

func (c *Context) EnqueueUpsertSubscription(appID uint, userID uint, planID uint) error {
	payload := UpsertSubscriptionPayload{
		AppID:  appID,
		UserID: userID,
		PlanID: planID,
	}

	if err := c.Queue().Enqueue("upsert_subscription", payload, "billing_queue"); err != nil {
		log.Printf("Error enqueueing upsert_subscription task: %v", err)
		return err
	}

	log.Printf("Enqueued upsert_subscription task for user ID %d", userID)

	return nil
}

func (c *Context) GetUpsertSubscriptionPayload() (*UpsertSubscriptionPayload, error) {
	var p UpsertSubscriptionPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}
