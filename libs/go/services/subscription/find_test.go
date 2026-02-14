package subscription_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/subscription"
	"gorm.io/gorm"
)

func TestFind_NoExisting_ReturnsCreate(t *testing.T) {
	db := newMockDB() // no subscriptions → not found
	svc := subscription.NewService(db)

	r, err := svc.Find(&subscription.FindInput{AppID: 1, UserID: 2, PlanID: 10})
	if err != nil {
		t.Fatal(err)
	}
	if r.Action != "create" {
		t.Fatalf("want create, got %s", r.Action)
	}
}

func TestFind_Existing_ReturnsUpdate(t *testing.T) {
	db := newMockDB()
	db.subscriptions[42] = &models.Subscription{
		Model: gorm.Model{ID: 42},
		AppID: 1, UserID: 2, Status: "active",
	}
	svc := subscription.NewService(db)

	r, err := svc.Find(&subscription.FindInput{AppID: 1, UserID: 2, PlanID: 10})
	if err != nil {
		t.Fatal(err)
	}
	if r.Action != "update" {
		t.Fatalf("want update, got %s", r.Action)
	}
}
