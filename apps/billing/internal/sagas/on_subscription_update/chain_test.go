package on_subscription_update_test

// Chain tests for the subscription update saga.
// These tests mock the DB and queue to verify that plan upgrades and
// downgrades enqueue the correct tasks end-to-end.

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_update"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

// --------------------------------------------------------------------------
// Test: Plan upgrade → prorated invoice + payment + entitlements
//
// Chain: update_subscription
//          → create_subscription_items (bulk)
//          → process_plan_switch (upgrade detected)
//              → create_upgrade_invoice → pay_invoice → resolve_invoice
//              → create_entitlements → create_single_entitlement
// --------------------------------------------------------------------------

func TestChain_PlanUpgrade_FullChain(t *testing.T) {
	t.Setenv("INVOICE_APP_URL", "https://test.example.com")

	addrID := uint(500)
	db := saga.NewStubDB()

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
	oldFixedItem := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 999, Quantity: 1}
	oldFixedItem.ID = 100
	db.Store(&oldFixedItem)

	// New plan fixed item ($29.99)
	newFixedItem := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 2999, Quantity: 1, PublicTitle: "Pro Plan"}
	newFixedItem.ID = 200
	db.Store(&newFixedItem)

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
		case *models.PlanItem:
			// FindFixedPlanItemForPlan — match plan_id from conditions
			for _, cond := range conds {
				if id, ok := cond.(uint); ok {
					if id == 10 {
						*d = oldFixedItem
						return nil
					}
					if id == 20 {
						*d = newFixedItem
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
		case *models.PaymentMethod:
			pm := models.PaymentMethod{AppID: 1, UserID: 2, ExternalID: "pm_123", ExternalType: "card"}
			pm.ID = 1
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
			pm := models.PaymentMethod{AppID: 1, UserID: 2, ExternalID: "pm_123", ExternalType: "card"}
			pm.ID = 1
			*methods = []models.PaymentMethod{pm}
		}
		return nil
	}

	runner := saga.NewRunner(db, nil,
		on_subscription_update.Steps,
		on_subscription_create.Steps,
		on_invoice_create.Steps,
		on_payment.Steps,
	)

	err := runner.Run("update_subscription", subscription.UpdateInput{
		SubscriptionID: 50,
		PlanID:         20,
		AppID:          1,
	})
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}

	// Verify the update handler executed
	if runner.Executed[0] != "update_subscription" {
		t.Fatalf("expected update_subscription first, got %q", runner.Executed[0])
	}

	if !runner.HasTask("pay_invoice") {
		t.Fatalf("expected pay_invoice\nfull: %v", runner.Executed)
	}
	if !runner.HasTask("resolve_invoice") {
		t.Fatalf("expected resolve_invoice\nfull: %v", runner.Executed)
	}
	payTasks := runner.TaskPayloads("pay_invoice")
	if len(payTasks) != 2 {
		t.Fatalf("expected two pay_invoice tasks (base + upgrade), got %d", len(payTasks))
	}

	entitlementCount := 0
	for _, created := range db.Created {
		if _, ok := created.(*models.Entitlement); ok {
			entitlementCount++
		}
	}
	if entitlementCount != 1 {
		t.Fatalf("expected one entitlement persisted, got %d", entitlementCount)
	}
}

// --------------------------------------------------------------------------
// Test: Plan downgrade → no upgrade invoice, entitlements synced immediately
//
// Chain: update_subscription
//          → create_subscription_items (bulk)
//          → process_plan_switch (downgrade detected)
//              → create_entitlements
// --------------------------------------------------------------------------

func TestChain_PlanDowngrade_NoInvoice(t *testing.T) {
	addrID := uint(500)
	db := saga.NewStubDB()

	oldPlan := models.Plan{Name: "Pro", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1, InvoiceDueByDays: 14}
	oldPlan.ID = 20
	db.Store(&oldPlan)

	newPlan := models.Plan{Name: "Basic", Interval: "month", IntervalCount: 1, Currency: "usd", AppID: 1, InvoiceDueByDays: 14}
	newPlan.ID = 10
	db.Store(&newPlan)

	user := models.User{Email: "jane@test.com", BillingAddressID: &addrID}
	user.ID = 2
	db.Store(&user)

	company := models.Company{Name: "TestCo", BillingAddressID: 500}
	company.ID = 1
	db.Store(&company)

	oldPlanID := uint(20)
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

	// Old plan fixed item ($29.99) → downgrading to $9.99
	oldFixedItem := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 2999, Quantity: 1}
	oldFixedItem.ID = 100
	db.Store(&oldFixedItem)

	newFixedItem := models.PlanItem{AppID: 1, PricingModel: "fixed", UnitAmount: 999, Quantity: 1, PublicTitle: "Basic Plan"}
	newFixedItem.ID = 200
	db.Store(&newFixedItem)

	db.ListFn = func(dest any, conds []any) error { return nil }

	db.FindFirstFn = func(dest any, conds []any) error {
		switch d := dest.(type) {
		case *models.PlanItem:
			for _, cond := range conds {
				if id, ok := cond.(uint); ok {
					if id == 20 {
						*d = oldFixedItem
						return nil
					}
					if id == 10 {
						*d = newFixedItem
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
		}
		return gorm.ErrRecordNotFound
	}

	db.ListIDsFn = func(table string, dest any, conds []any) error {
		ids := dest.(*[]uint)
		switch table {
		case "plan_items":
			*ids = []uint{200} // new plan items
		case "subscription_items":
			var created []uint
			for _, c := range db.Created {
				if si, ok := c.(*models.SubscriptionItem); ok {
					created = append(created, si.ID)
				}
			}
			*ids = created
		}
		return nil
	}

	runner := saga.NewRunner(db, nil,
		on_subscription_update.Steps,
		on_subscription_create.Steps,
	)

	err := runner.Run("update_subscription", subscription.UpdateInput{
		SubscriptionID: 50,
		PlanID:         10,
		AppID:          1,
	})
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}

	// Verify a downgrade does NOT create an upgrade invoice
	if runner.HasTask("create_upgrade_invoice") {
		t.Fatal("downgrade should not create an upgrade invoice")
	}

	if !runner.HasTask("pay_invoice") {
		t.Fatalf("expected pay_invoice for normal renewal invoice\nfull: %v", runner.Executed)
	}
}
