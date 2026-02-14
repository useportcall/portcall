package subscription_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

func TestCreate_Success(t *testing.T) {
	db := newMockDB()
	planID := uint(10)
	db.plans[planID] = &models.Plan{
		Model: gorm.Model{ID: planID}, Name: "Pro",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", AppID: 1,
	}
	addrID := uint(100)
	db.users[2] = &models.User{
		Model: gorm.Model{ID: 2}, BillingAddressID: &addrID,
	}

	svc := subscription.NewService(db)
	r, err := svc.Create(&subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: planID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Subscription.Status != "active" {
		t.Fatalf("want active, got %s", r.Subscription.Status)
	}
	if r.Subscription.UserID != 2 {
		t.Fatalf("want user 2, got %d", r.Subscription.UserID)
	}
	if r.PlanID != planID {
		t.Fatalf("want plan %d, got %d", planID, r.PlanID)
	}
}

func TestCreate_PaidPlan_SetsRollback(t *testing.T) {
	db := newMockDB()
	planID := uint(10)
	freePlanID := uint(5)
	db.plans[planID] = &models.Plan{
		Model: gorm.Model{ID: planID}, Name: "Pro",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", AppID: 1, IsFree: false,
	}
	db.users[2] = &models.User{Model: gorm.Model{ID: 2}}
	db.freePlans = []models.Plan{
		{Model: gorm.Model{ID: freePlanID}, IsFree: true},
	}

	svc := subscription.NewService(db)
	r, err := svc.Create(&subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: planID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Subscription.RollbackPlanID == nil {
		t.Fatal("expected rollback plan set")
	}
	if *r.Subscription.RollbackPlanID != freePlanID {
		t.Fatalf("want rollback %d, got %d",
			freePlanID, *r.Subscription.RollbackPlanID)
	}
}

func TestCreate_FreePlan_NoRollback(t *testing.T) {
	db := newMockDB()
	planID := uint(5)
	db.plans[planID] = &models.Plan{
		Model: gorm.Model{ID: planID}, Name: "Free",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", AppID: 1, IsFree: true,
	}
	db.users[2] = &models.User{Model: gorm.Model{ID: 2}}

	svc := subscription.NewService(db)
	r, err := svc.Create(&subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: planID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Subscription.RollbackPlanID != nil {
		t.Fatal("expected no rollback plan for free plan")
	}
}

func TestCreate_PaidPlan_FallbacksToCheapestFixedPlan(t *testing.T) {
	db := newMockDB()
	currentID := uint(10)
	fallbackID := uint(20)
	db.plans[currentID] = &models.Plan{
		Model: gorm.Model{ID: currentID}, Name: "Pro",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", AppID: 1, IsFree: false,
	}
	db.plans[fallbackID] = &models.Plan{
		Model: gorm.Model{ID: fallbackID}, Name: "Starter",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", AppID: 1, IsFree: false,
	}
	db.users[2] = &models.User{Model: gorm.Model{ID: 2}}
	db.planItemsByPlan[fallbackID] = &models.PlanItem{
		Model: gorm.Model{ID: 100}, PlanID: fallbackID,
		PricingModel: "fixed", UnitAmount: 1000,
	}

	svc := subscription.NewService(db)
	r, err := svc.Create(&subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: currentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Subscription.RollbackPlanID == nil || *r.Subscription.RollbackPlanID != fallbackID {
		t.Fatalf("expected rollback %d, got %v", fallbackID, r.Subscription.RollbackPlanID)
	}
}

func TestCreate_CreatesSubscriptionItemsTransactionally(t *testing.T) {
	db := newMockDB()
	planID := uint(10)
	itemID := uint(50)
	db.plans[planID] = &models.Plan{
		Model: gorm.Model{ID: planID}, Name: "Pro",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", AppID: 1,
	}
	db.users[2] = &models.User{Model: gorm.Model{ID: 2}}
	db.planItemIDs = []uint{itemID}
	db.planItems[itemID] = &models.PlanItem{
		Model: gorm.Model{ID: itemID}, PlanID: planID, AppID: 1,
		Quantity: 1, UnitAmount: 2999, PricingModel: "fixed",
		Interval: "month", IntervalCount: 1,
	}

	svc := subscription.NewService(db)
	r, err := svc.Create(&subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: planID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.ItemCount != 1 {
		t.Fatalf("expected 1 item, got %d", r.ItemCount)
	}
	if len(db.subItems) != 1 {
		t.Fatalf("expected 1 persisted subscription item, got %d", len(db.subItems))
	}
}
