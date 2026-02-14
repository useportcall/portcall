package on_subscription_create_test

// Chain tests for the subscription creation saga.
// These tests mock the DB and queue to verify that the correct
// tasks are enqueued with the right data end-to-end.

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

// --------------------------------------------------------------------------
// Test: Full subscription creation with 2 items + 2 features
//
// Chain: create_subscription (subscription + items in one txn)
//          → create_invoice → pay_invoice → resolve_invoice → send_invoice_paid_email
//          → create_entitlements → create_single_entitlement (×2)
// --------------------------------------------------------------------------

func TestChain_CreateSubscription_FullChain(t *testing.T) {
	t.Setenv("INVOICE_APP_URL", "https://test.example.com")

	addrID := uint(500)
	db := saga.NewStubDB()

	plan := models.Plan{Name: "Pro", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1, InvoiceDueByDays: 14}
	plan.ID = 10
	db.Store(&plan)

	user := models.User{Email: "user@test.com", Name: "Jane", BillingAddressID: &addrID, PaymentCustomerID: "cus_1"}
	user.ID = 2
	db.Store(&user)

	company := models.Company{Name: "TestCo", BillingAddressID: 500}
	company.ID = 1
	db.Store(&company)

	planItem1 := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 1999, Quantity: 1, PublicTitle: "Base Fee"}
	planItem1.ID = 100
	db.Store(&planItem1)

	planItem2 := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 500, Quantity: 1, PublicTitle: "Add-on"}
	planItem2.ID = 101
	db.Store(&planItem2)

	feat1 := models.Feature{PublicID: "feat_api", IsMetered: true}
	feat1.ID = 200
	db.Store(&feat1)

	feat2 := models.Feature{PublicID: "feat_storage", IsMetered: false}
	feat2.ID = 201
	db.Store(&feat2)

	pf1 := models.PlanFeature{AppID: 1, FeatureID: 200, Interval: "month", Quota: 10000}
	pf1.ID = 300
	db.Store(&pf1)

	pf2 := models.PlanFeature{AppID: 1, FeatureID: 201, Interval: "month", Quota: 50}
	pf2.ID = 301
	db.Store(&pf2)

	connection := models.Connection{Source: "local", PublicKey: "pk_test"}
	connection.ID = 1
	db.Store(&connection)

	pm := models.PaymentMethod{AppID: 1, UserID: 2, ExternalID: "pm_123", ExternalType: "card"}
	pm.ID = 1
	db.Store(&pm)

	db.ListFn = func(dest any, conds []any) error {
		if plans, ok := dest.(*[]models.Plan); ok {
			// FindFreePlanForApp — no free plan
			*plans = nil
			return nil
		}
		return nil
	}

	db.FindFirstFn = func(dest any, conds []any) error {
		switch d := dest.(type) {
		case *models.Company:
			*d = company
			return nil
		case *models.Invoice:
			return gorm.ErrRecordNotFound
		case *models.Entitlement:
			return gorm.ErrRecordNotFound
		case *models.PaymentMethod:
			*d = pm
			return nil
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
			*ids = []uint{300, 301}
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
			*methods = []models.PaymentMethod{pm}
		}
		return nil
	}

	runner := saga.NewRunner(db, nil,
		on_subscription_create.Steps,
		on_invoice_create.Steps,
		on_payment.Steps,
	)

	err := runner.Run("create_subscription", subscription.CreateInput{
		AppID:  1,
		UserID: 2,
		PlanID: 10,
	})
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}

	if len(runner.Executed) == 0 || runner.Executed[0] != "create_subscription" {
		t.Fatalf("expected create_subscription first, got %v", runner.Executed)
	}

	// Verify the chain continues through payment
	if !runner.HasTask("pay_invoice") {
		t.Fatalf("expected pay_invoice in chain\nfull: %v", runner.Executed)
	}

	if !runner.HasTask("resolve_invoice") {
		t.Fatalf("expected resolve_invoice in chain\nfull: %v", runner.Executed)
	}

	// Verify email is enqueued (on email queue, not executed)
	if !runner.HasTask("send_invoice_paid_email") {
		t.Fatalf("expected send_invoice_paid_email in chain\nfull: %v", runner.Executed)
	}

	entitlementCount := 0
	for _, created := range db.Created {
		if _, ok := created.(*models.Entitlement); ok {
			entitlementCount++
		}
	}
	if entitlementCount != 2 {
		t.Fatalf("expected 2 entitlements persisted, got %d", entitlementCount)
	}
}

// --------------------------------------------------------------------------
// Test: Subscription with no plan items → item chain skipped
//
// Chain: create_subscription
//          → create_entitlements → create_single_entitlement
// --------------------------------------------------------------------------

func TestChain_CreateSubscription_NoItems(t *testing.T) {
	addrID := uint(500)
	db := saga.NewStubDB()

	plan := models.Plan{Name: "Basic", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1}
	plan.ID = 10
	db.Store(&plan)

	user := models.User{Email: "u@test.com", BillingAddressID: &addrID}
	user.ID = 2
	db.Store(&user)

	feat := models.Feature{PublicID: "feat_basic", IsMetered: false}
	feat.ID = 200
	db.Store(&feat)

	pf := models.PlanFeature{AppID: 1, FeatureID: 200, Interval: "month", Quota: 100}
	pf.ID = 300
	db.Store(&pf)

	db.ListFn = func(dest any, conds []any) error { return nil }

	db.FindFirstFn = func(dest any, conds []any) error {
		if _, ok := dest.(*models.Entitlement); ok {
			return gorm.ErrRecordNotFound
		}
		return gorm.ErrRecordNotFound
	}

	db.ListIDsFn = func(table string, dest any, conds []any) error {
		ids := dest.(*[]uint)
		switch table {
		case "plan_items":
			*ids = []uint{}
		case "plan_features":
			*ids = []uint{300}
		case "entitlements":
			*ids = []uint{}
		}
		return nil
	}

	runner := saga.NewRunner(db, nil, on_subscription_create.Steps)

	err := runner.Run("create_subscription", subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: 10,
	})
	if err != nil {
		t.Fatalf("chain failed: %v", err)
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
	if runner.HasTask("pay_invoice") {
		t.Fatal("should not enqueue payment when plan has no items")
	}
}

// --------------------------------------------------------------------------
// Test: Subscription with no features → entitlement chain skipped
//
// Chain: create_subscription
//          → create_invoice
//          → create_entitlements (no features → stops)
// --------------------------------------------------------------------------

func TestChain_CreateSubscription_NoFeatures(t *testing.T) {
	addrID := uint(500)
	db := saga.NewStubDB()

	plan := models.Plan{Name: "Basic", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1}
	plan.ID = 10
	db.Store(&plan)

	user := models.User{Email: "u@test.com", BillingAddressID: &addrID}
	user.ID = 2
	db.Store(&user)

	planItem := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 999, Quantity: 1}
	planItem.ID = 100
	db.Store(&planItem)

	company := models.Company{Name: "TestCo", BillingAddressID: 500}
	company.ID = 1
	db.Store(&company)

	db.ListFn = func(dest any, conds []any) error { return nil }

	db.FindFirstFn = func(dest any, conds []any) error {
		if d, ok := dest.(*models.Company); ok {
			*d = company
			return nil
		}
		return gorm.ErrRecordNotFound
	}

	db.ListIDsFn = func(table string, dest any, conds []any) error {
		ids := dest.(*[]uint)
		switch table {
		case "plan_items":
			*ids = []uint{100}
		case "plan_features":
			*ids = []uint{}
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

	runner := saga.NewRunner(db, nil, on_subscription_create.Steps)

	err := runner.Run("create_subscription", subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: 10,
	})
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}

	entitlementCount := 0
	for _, created := range db.Created {
		if _, ok := created.(*models.Entitlement); ok {
			entitlementCount++
		}
	}
	if entitlementCount != 0 {
		t.Fatalf("expected no entitlements persisted, got %d", entitlementCount)
	}
	if !runner.HasTask("pay_invoice") {
		t.Fatalf("expected pay_invoice\nfull: %v", runner.Executed)
	}
}
