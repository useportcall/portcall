package invoice_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	inv "github.com/useportcall/portcall/libs/go/services/invoice"
)

func TestList_AllInvoices(t *testing.T) {
	db := &mockDB{
		invoices: []models.Invoice{
			{PublicID: "inv_1"}, {PublicID: "inv_2"},
		},
	}

	svc := inv.NewService(db)
	result, err := svc.List(&inv.ListInput{AppID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Invoices) != 2 {
		t.Fatalf("expected 2 invoices, got %d", len(result.Invoices))
	}
}

func TestList_BySubscription(t *testing.T) {
	sub := &models.Subscription{}
	sub.ID = 10

	db := &mockDB{
		subscriptions: map[string]*models.Subscription{"sub_abc": sub},
		invoices:      []models.Invoice{{PublicID: "inv_1"}},
	}

	svc := inv.NewService(db)
	result, err := svc.List(&inv.ListInput{AppID: 1, SubscriptionID: "sub_abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Invoices) != 1 {
		t.Fatalf("expected 1 invoice, got %d", len(result.Invoices))
	}
}

func TestList_SubscriptionNotFound(t *testing.T) {
	db := &mockDB{
		subscriptions: map[string]*models.Subscription{},
	}

	svc := inv.NewService(db)
	_, err := svc.List(&inv.ListInput{AppID: 1, SubscriptionID: "sub_missing"})
	if err == nil {
		t.Fatal("expected error for missing subscription")
	}
}

func TestList_ByUser(t *testing.T) {
	user := &models.User{}
	user.ID = 5

	db := &mockDB{
		users:    map[string]*models.User{"usr_xyz": user},
		invoices: []models.Invoice{{PublicID: "inv_1"}},
	}

	svc := inv.NewService(db)
	result, err := svc.List(&inv.ListInput{AppID: 1, UserID: "usr_xyz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Invoices) != 1 {
		t.Fatalf("expected 1 invoice, got %d", len(result.Invoices))
	}
}
