package subscription_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

func TestCreateItems_Success(t *testing.T) {
	db := newMockDB()
	db.planItemIDs = []uint{100, 200}
	db.plans[10] = &models.Plan{
		Model: gorm.Model{ID: 10}, Name: "Pro",
	}
	db.planItems[100] = &models.PlanItem{
		Model: gorm.Model{ID: 100}, PricingModel: "fixed",
		AppID: 1, UnitAmount: 1000,
	}
	db.planItems[200] = &models.PlanItem{
		Model: gorm.Model{ID: 200}, PricingModel: "metered",
		AppID: 1, UnitAmount: 50, PublicTitle: "API Calls",
	}

	svc := subscription.NewService(db)
	r, err := svc.CreateItems(&subscription.CreateItemsInput{
		SubscriptionID: 1, PlanID: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.ItemCount != 2 {
		t.Fatalf("want 2 items, got %d", r.ItemCount)
	}
	if len(db.created) != 2 {
		t.Fatalf("want 2 creates, got %d", len(db.created))
	}
}

func TestCreateItems_NoItems(t *testing.T) {
	db := newMockDB()
	db.planItemIDs = []uint{}

	svc := subscription.NewService(db)
	r, err := svc.CreateItems(&subscription.CreateItemsInput{
		SubscriptionID: 1, PlanID: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.ItemCount != 0 {
		t.Fatalf("want 0 items, got %d", r.ItemCount)
	}
}
