package subscription_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

func TestUpdate_Success(t *testing.T) {
	db := newMockDB()
	oldPlanID := uint(10)
	newPlanID := uint(20)

	db.subscriptions[1] = &models.Subscription{
		Model: gorm.Model{ID: 1}, PlanID: &oldPlanID,
		Currency: "USD", BillingInterval: "month",
		BillingIntervalCount: 1,
	}
	db.plans[oldPlanID] = &models.Plan{
		Model: gorm.Model{ID: oldPlanID},
		Currency: "USD", Interval: "month", IntervalCount: 1,
		InvoiceDueByDays: 10,
	}
	db.plans[newPlanID] = &models.Plan{
		Model: gorm.Model{ID: newPlanID},
		Currency: "USD", Interval: "month", IntervalCount: 1,
		InvoiceDueByDays: 14,
	}

	svc := subscription.NewService(db)
	r, err := svc.Update(&subscription.UpdateInput{
		SubscriptionID: 1, PlanID: newPlanID, AppID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.OldPlanID != oldPlanID || r.NewPlanID != newPlanID {
		t.Fatal("plan IDs mismatch")
	}
	if len(db.deleted) == 0 {
		t.Fatal("expected old items deleted")
	}
}

func TestUpdate_CurrencyMismatch(t *testing.T) {
	db := newMockDB()
	oldPlanID := uint(10)

	db.subscriptions[1] = &models.Subscription{
		Model: gorm.Model{ID: 1}, PlanID: &oldPlanID,
		Currency: "USD", BillingInterval: "month",
		BillingIntervalCount: 1,
	}
	db.plans[oldPlanID] = &models.Plan{
		Model: gorm.Model{ID: oldPlanID},
		Currency: "USD", Interval: "month", IntervalCount: 1,
	}
	db.plans[20] = &models.Plan{
		Model: gorm.Model{ID: 20},
		Currency: "EUR", Interval: "month", IntervalCount: 1,
	}

	svc := subscription.NewService(db)
	_, err := svc.Update(&subscription.UpdateInput{
		SubscriptionID: 1, PlanID: 20,
	})
	if err == nil {
		t.Fatal("expected currency mismatch error")
	}
}

func TestUpdate_IntervalMismatch(t *testing.T) {
	db := newMockDB()
	oldPlanID := uint(10)

	db.subscriptions[1] = &models.Subscription{
		Model: gorm.Model{ID: 1}, PlanID: &oldPlanID,
		Currency: "USD", BillingInterval: "month",
		BillingIntervalCount: 1,
	}
	db.plans[oldPlanID] = &models.Plan{
		Model: gorm.Model{ID: oldPlanID},
		Currency: "USD", Interval: "month", IntervalCount: 1,
	}
	db.plans[20] = &models.Plan{
		Model: gorm.Model{ID: 20},
		Currency: "USD", Interval: "year", IntervalCount: 1,
	}

	svc := subscription.NewService(db)
	_, err := svc.Update(&subscription.UpdateInput{
		SubscriptionID: 1, PlanID: 20,
	})
	if err == nil {
		t.Fatal("expected interval mismatch error")
	}
}
