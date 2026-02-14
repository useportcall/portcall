//go:build integration

package subscription_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/services/subscription"
)

func TestCreate_PaidPlanUsesFreeRollbackPlan(t *testing.T) {
	app, user := seedIntegrationAppUser(t)
	free := seedPlan(t, app.ID, "Free", true)
	paid := seedPlan(t, app.ID, "Pro", false)

	svc := subscription.NewService(integrationDB)
	result, err := svc.Create(&subscription.CreateInput{
		AppID: app.ID, UserID: user.ID, PlanID: paid.ID,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.Subscription.RollbackPlanID == nil {
		t.Fatal("expected rollback plan to be set")
	}
	if *result.Subscription.RollbackPlanID != free.ID {
		t.Fatalf(
			"expected rollback plan %d, got %d",
			free.ID,
			*result.Subscription.RollbackPlanID,
		)
	}
}
