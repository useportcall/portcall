package subscription_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

func TestCreate_RollsBackIfPlanItemLookupFails(t *testing.T) {
	db := newMockDB()
	planID := uint(10)
	db.plans[planID] = &models.Plan{
		Model: gorm.Model{ID: planID}, Name: "Pro",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", AppID: 1,
	}
	db.users[2] = &models.User{Model: gorm.Model{ID: 2}}
	db.planItemIDs = []uint{999}

	svc := subscription.NewService(db)
	_, err := svc.Create(&subscription.CreateInput{
		AppID: 1, UserID: 2, PlanID: planID,
	})
	if err == nil {
		t.Fatal("expected error for missing plan item")
	}
	if len(db.subscriptions) != 0 {
		t.Fatalf("expected transaction rollback, got %d subscriptions", len(db.subscriptions))
	}
	if len(db.subItems) != 0 {
		t.Fatalf("expected no items after rollback, got %d", len(db.subItems))
	}
}
