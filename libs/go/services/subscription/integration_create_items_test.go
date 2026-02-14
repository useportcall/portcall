//go:build integration

package subscription_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func TestCreate_PersistsSubscriptionAndItems(t *testing.T) {
	app, user := seedIntegrationAppUser(t)
	plan := seedPlan(t, app.ID, "Pro", false)
	seedPlanItem(t, app.ID, plan.ID, "Base", 2900)
	seedPlanItem(t, app.ID, plan.ID, "Addon", 900)

	svc := subscription.NewService(integrationDB)
	result, err := svc.Create(&subscription.CreateInput{
		AppID: app.ID, UserID: user.ID, PlanID: plan.ID,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.ItemCount != 2 {
		t.Fatalf("expected ItemCount=2, got %d", result.ItemCount)
	}

	var count int64
	if err := integrationDB.Count(
		&count,
		&models.SubscriptionItem{},
		"subscription_id = ?",
		result.Subscription.ID,
	); err != nil {
		t.Fatalf("count subscription items: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 persisted subscription items, got %d", count)
	}
}
