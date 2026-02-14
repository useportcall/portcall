package subscription_test

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

func TestPlanSwitch_Upgrade(t *testing.T) {
	db := newMockDB()
	db.planItemsByPlan[10] = &models.PlanItem{
		Model: gorm.Model{ID: 1}, UnitAmount: 1000,
	}
	db.planItemsByPlan[20] = &models.PlanItem{
		Model: gorm.Model{ID: 2}, UnitAmount: 3000,
	}
	db.subscriptions[1] = &models.Subscription{
		Model: gorm.Model{ID: 1}, UserID: 2,
	}

	svc := subscription.NewService(db)
	r, err := svc.PlanSwitch(&subscription.PlanSwitchInput{
		OldPlanID: 10, NewPlanID: 20, SubscriptionID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !r.IsUpgrade {
		t.Fatal("expected upgrade")
	}
	if r.PriceDifference != 2000 {
		t.Fatalf("want 2000, got %d", r.PriceDifference)
	}
}

func TestPlanSwitch_Downgrade(t *testing.T) {
	db := newMockDB()
	db.planItemsByPlan[10] = &models.PlanItem{
		Model: gorm.Model{ID: 1}, UnitAmount: 3000,
	}
	db.planItemsByPlan[20] = &models.PlanItem{
		Model: gorm.Model{ID: 2}, UnitAmount: 1000,
	}

	svc := subscription.NewService(db)
	r, err := svc.PlanSwitch(&subscription.PlanSwitchInput{
		OldPlanID: 10, NewPlanID: 20, SubscriptionID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.IsUpgrade {
		t.Fatal("expected downgrade")
	}
}

func TestPlanSwitch_MidPeriod(t *testing.T) {
	now := time.Now()
	db := newMockDB()
	db.planItemsByPlan[10] = &models.PlanItem{
		Model: gorm.Model{ID: 1}, UnitAmount: 1000,
	}
	db.planItemsByPlan[20] = &models.PlanItem{
		Model: gorm.Model{ID: 2}, UnitAmount: 3000,
	}
	db.subscriptions[1] = &models.Subscription{
		Model:       gorm.Model{ID: 1},
		UserID:      2,
		LastResetAt: now.AddDate(0, 0, -15),
		NextResetAt: now.AddDate(0, 0, 15),
	}

	svc := subscription.NewService(db)
	r, err := svc.PlanSwitch(&subscription.PlanSwitchInput{
		OldPlanID: 10, NewPlanID: 20, SubscriptionID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !r.IsUpgrade {
		t.Fatal("expected upgrade")
	}
	if r.PriceDifference < 900 || r.PriceDifference > 1100 {
		t.Fatalf("expected prorated ~1000, got %d", r.PriceDifference)
	}
}

func TestProrate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name   string
		amount int64
		start  time.Time
		end    time.Time
		want   int64
	}{
		{"full", 1000, now, now.AddDate(0, 1, 0), 1000},
		{"half", 1000, now.Add(-12 * time.Hour), now.Add(12 * time.Hour), 500},
		{"zero remaining", 1000, now.Add(-24 * time.Hour), now, 1000},
		{"invalid", 1000, now, now, 1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subscription.Prorate(tt.amount, tt.start, tt.end, now)
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
