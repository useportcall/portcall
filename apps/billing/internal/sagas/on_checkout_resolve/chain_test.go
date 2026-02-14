package on_checkout_resolve_test

// Chain tests for the checkout resolution saga.
// These tests mock the DB and queue to verify that the correct
// tasks are enqueued with the right data when a checkout completes.
//
// The chain starts at resolve_checkout_session (skipping the Stripe
// webhook adapter which is tested separately in unit tests).

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_checkout_resolve"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_update"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
	"gorm.io/gorm"
)

// --------------------------------------------------------------------------
// Test: New customer checkout → create subscription → items → invoice → payment
//
// Chain: resolve_checkout_session → create_payment_method → upsert_subscription
//          → create_subscription (subscription + items in one txn)
//              → create_invoice → pay_invoice → resolve_invoice → send_invoice_paid_email
//          → create_entitlements → create_single_entitlement (×1)
// --------------------------------------------------------------------------

func TestChain_NewCustomerCheckout(t *testing.T) {
	t.Setenv("INVOICE_APP_URL", "https://test.example.com")

	addrID := uint(500)
	db := saga.NewStubDB()

	// Checkout session to resolve
	session := models.CheckoutSession{
		AppID:             1,
		PlanID:            10,
		UserID:            2,
		ExternalSessionID: "seti_123",
		Status:            "active",
	}
	session.ID = 1
	db.Store(&session)

	plan := models.Plan{Name: "Pro", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1, InvoiceDueByDays: 14}
	plan.ID = 10
	db.Store(&plan)

	user := models.User{Email: "jane@test.com", Name: "Jane", BillingAddressID: &addrID, PaymentCustomerID: "cus_1"}
	user.ID = 2
	db.Store(&user)

	company := models.Company{Name: "TestCo", BillingAddressID: 500}
	company.ID = 1
	db.Store(&company)

	planItem1 := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 2999, Quantity: 1, PublicTitle: "Pro Plan"}
	planItem1.ID = 100
	db.Store(&planItem1)

	planItem2 := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 500, Quantity: 1, PublicTitle: "Add-on"}
	planItem2.ID = 101
	db.Store(&planItem2)

	feat := models.Feature{PublicID: "feat_api", IsMetered: true}
	feat.ID = 200
	db.Store(&feat)

	pf := models.PlanFeature{AppID: 1, FeatureID: 200, Interval: "month", Quota: 10000}
	pf.ID = 300
	db.Store(&pf)

	connection := models.Connection{Source: "local", PublicKey: "pk_test"}
	connection.ID = 1
	db.Store(&connection)

	db.ListFn = func(dest any, conds []any) error { return nil }

	db.FindFirstFn = func(dest any, conds []any) error {
		switch d := dest.(type) {
		case *models.CheckoutSession:
			*d = session
			return nil
		case *models.PaymentMethod:
			// First call: UpsertPaymentMethod check → not found → create
			// Subsequent calls: FindPaymentMethodForUser → return created
			for _, c := range db.Created {
				if pm, ok := c.(*models.PaymentMethod); ok {
					*d = *pm
					return nil
				}
			}
			return gorm.ErrRecordNotFound
		case *models.Subscription:
			// No existing subscription → create flow
			return gorm.ErrRecordNotFound
		case *models.Company:
			*d = company
			return nil
		case *models.Invoice:
			return gorm.ErrRecordNotFound
		case *models.Entitlement:
			return gorm.ErrRecordNotFound
		case *models.Connection:
			*d = connection
			return nil
		}
		return gorm.ErrRecordNotFound
	}

	db.ListIDsFn = func(table string, dest any, conds []any) error {
		ids := dest.(*[]uint)
		switch table {
		case "plan_items":
			*ids = []uint{100, 101}
		case "plan_features":
			*ids = []uint{300}
		case "subscription_items":
			var created []uint
			for _, c := range db.Created {
				if si, ok := c.(*models.SubscriptionItem); ok {
					created = append(created, si.ID)
				}
			}
			*ids = created
		case "entitlements":
			*ids = []uint{}
		}
		return nil
	}

	db.CountFn = func(count *int64, dest any, query string, args []any) error {
		*count = 0
		return nil
	}

	db.ListWithOrderAndLimitFn = func(dest any, order string, limit int, conds []any) error {
		if methods, ok := dest.(*[]models.PaymentMethod); ok {
			for _, c := range db.Created {
				if pm, ok := c.(*models.PaymentMethod); ok {
					*methods = append(*methods, *pm)
				}
			}
		}
		return nil
	}

	runner := saga.NewRunner(db, nil,
		on_checkout_resolve.Steps,
		on_subscription_create.Steps,
		on_subscription_update.Steps,
		on_invoice_create.Steps,
		on_payment.Steps,
	)

	err := runner.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       "seti_123",
		ExternalPaymentMethodID: "pm_external_456",
	})
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}

	// Verify the full chain executed in order
	if !runner.HasTask("resolve_checkout_session") {
		t.Fatal("missing resolve_checkout_session")
	}
	if !runner.HasTask("create_payment_method") {
		t.Fatalf("missing create_payment_method\nfull: %v", runner.Executed)
	}
	if !runner.HasTask("upsert_subscription") {
		t.Fatalf("missing upsert_subscription\nfull: %v", runner.Executed)
	}
	if !runner.HasTask("create_subscription") {
		t.Fatalf("expected create (not update) for new customer\nfull: %v", runner.Executed)
	}
	if runner.HasTask("update_subscription") {
		t.Fatal("should not update when no existing subscription")
	}

	if !runner.HasTask("pay_invoice") {
		t.Fatalf("expected pay_invoice\nfull: %v", runner.Executed)
	}
	if !runner.HasTask("resolve_invoice") {
		t.Fatalf("expected resolve_invoice\nfull: %v", runner.Executed)
	}
	if !runner.HasTask("send_invoice_paid_email") {
		t.Fatalf("expected send_invoice_paid_email\nfull: %v", runner.Executed)
	}

	entitlementCount := 0
	for _, created := range db.Created {
		if _, ok := created.(*models.Entitlement); ok {
			entitlementCount++
		}
	}
	if entitlementCount != 1 {
		t.Fatalf("expected 1 entitlement persisted, got %d", entitlementCount)
	}
}

// --------------------------------------------------------------------------
// Test: Existing customer checkout → update subscription → upgrade invoice
//
// Chain: resolve_checkout_session → create_payment_method → upsert_subscription
//          → update_subscription
//              → create_subscription_items (bulk)
//              → process_plan_switch → create_upgrade_invoice → pay_invoice
//                                    → create_entitlements → create_single_entitlement
// --------------------------------------------------------------------------

func TestChain_ExistingCustomerCheckout_Upgrade(t *testing.T) {
	t.Setenv("INVOICE_APP_URL", "https://test.example.com")

	addrID := uint(500)
	db := saga.NewStubDB()

	session := models.CheckoutSession{
		AppID:             1,
		PlanID:            20,
		UserID:            2,
		ExternalSessionID: "seti_999",
		Status:            "active",
	}
	session.ID = 1
	db.Store(&session)

	oldPlan := models.Plan{Name: "Basic", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1, InvoiceDueByDays: 14}
	oldPlan.ID = 10
	db.Store(&oldPlan)

	newPlan := models.Plan{Name: "Pro", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1, InvoiceDueByDays: 14}
	newPlan.ID = 20
	db.Store(&newPlan)

	user := models.User{Email: "jane@test.com", Name: "Jane", BillingAddressID: &addrID, PaymentCustomerID: "cus_1"}
	user.ID = 2
	db.Store(&user)

	company := models.Company{Name: "TestCo", BillingAddressID: 500}
	company.ID = 1
	db.Store(&company)

	oldPlanID := uint(10)
	existingSub := models.Subscription{
		AppID:                1,
		UserID:               2,
		Status:               "active",
		Currency:             "usd",
		PlanID:               &oldPlanID,
		BillingInterval:      "month",
		BillingIntervalCount: 1,
		BillingAddressID:     &addrID,
		InvoiceDueByDays:     14,
		LastResetAt:          time.Now().AddDate(0, 0, -15),
		NextResetAt:          time.Now().AddDate(0, 0, 15),
	}
	existingSub.ID = 50
	db.Store(&existingSub)

	// Old plan fixed item ($9.99)
	oldItem := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 999, Quantity: 1}
	oldItem.ID = 100
	db.Store(&oldItem)

	// New plan fixed item ($29.99)
	newItem := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 2999, Quantity: 1, PublicTitle: "Pro Plan"}
	newItem.ID = 200
	db.Store(&newItem)

	feat := models.Feature{PublicID: "feat_pro", IsMetered: false}
	feat.ID = 300
	db.Store(&feat)

	pf := models.PlanFeature{AppID: 1, FeatureID: 300, Interval: "month", Quota: 50000}
	pf.ID = 400
	db.Store(&pf)

	connection := models.Connection{Source: "local", PublicKey: "pk_test"}
	connection.ID = 1
	db.Store(&connection)

	db.ListFn = func(dest any, conds []any) error { return nil }

	db.FindFirstFn = func(dest any, conds []any) error {
		switch d := dest.(type) {
		case *models.CheckoutSession:
			*d = session
			return nil
		case *models.PaymentMethod:
			for _, c := range db.Created {
				if pm, ok := c.(*models.PaymentMethod); ok {
					*d = *pm
					return nil
				}
			}
			return gorm.ErrRecordNotFound
		case *models.Subscription:
			// Existing subscription → update flow
			*d = existingSub
			return nil
		case *models.PlanItem:
			// FindFixedPlanItemForPlan — match by plan_id from conditions
			for _, cond := range conds {
				if id, ok := cond.(uint); ok {
					if id == 10 {
						*d = oldItem
						return nil
					}
					if id == 20 {
						*d = newItem
						return nil
					}
				}
			}
			return gorm.ErrRecordNotFound
		case *models.Company:
			*d = company
			return nil
		case *models.Invoice:
			return gorm.ErrRecordNotFound
		case *models.Entitlement:
			return gorm.ErrRecordNotFound
		case *models.Connection:
			*d = connection
			return nil
		}
		return gorm.ErrRecordNotFound
	}

	db.ListIDsFn = func(table string, dest any, conds []any) error {
		ids := dest.(*[]uint)
		switch table {
		case "plan_items":
			*ids = []uint{200} // new plan items
		case "plan_features":
			*ids = []uint{400}
		case "subscription_items":
			var created []uint
			for _, c := range db.Created {
				if si, ok := c.(*models.SubscriptionItem); ok {
					created = append(created, si.ID)
				}
			}
			*ids = created
		case "entitlements":
			*ids = []uint{}
		}
		return nil
	}

	db.CountFn = func(count *int64, dest any, query string, args []any) error {
		*count = 0
		return nil
	}

	db.ListWithOrderAndLimitFn = func(dest any, order string, limit int, conds []any) error {
		if methods, ok := dest.(*[]models.PaymentMethod); ok {
			for _, c := range db.Created {
				if pm, ok := c.(*models.PaymentMethod); ok {
					*methods = append(*methods, *pm)
				}
			}
		}
		return nil
	}

	runner := saga.NewRunner(db, nil,
		on_checkout_resolve.Steps,
		on_subscription_create.Steps,
		on_subscription_update.Steps,
		on_invoice_create.Steps,
		on_payment.Steps,
	)

	err := runner.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       "seti_999",
		ExternalPaymentMethodID: "pm_external_789",
	})
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}

	// Should route through update, not create
	if !runner.HasTask("update_subscription") {
		t.Fatalf("expected update_subscription for existing customer\nfull: %v", runner.Executed)
	}
	if runner.HasTask("create_subscription") {
		t.Fatal("should not create when subscription already exists")
	}

	entitlementCount := 0
	for _, created := range db.Created {
		if _, ok := created.(*models.Entitlement); ok {
			entitlementCount++
		}
	}
	if entitlementCount != 1 {
		t.Fatalf("expected 1 entitlement persisted, got %d", entitlementCount)
	}

	if !runner.HasTask("pay_invoice") {
		t.Fatalf("expected pay_invoice\nfull: %v", runner.Executed)
	}
	payTasks := runner.TaskPayloads("pay_invoice")
	if len(payTasks) != 2 {
		t.Fatalf("expected two pay_invoice tasks (base + upgrade), got %d", len(payTasks))
	}
}
