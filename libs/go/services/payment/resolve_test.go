package payment_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	pay "github.com/useportcall/portcall/libs/go/services/payment"
)

func TestResolve_Success(t *testing.T) {
	db := &mockDB{
		invoice: &models.Invoice{Total: 5000, UserID: 1, AppID: 1, InvoiceNumber: "INV-001", Currency: "USD"},
		user:    &models.User{Email: "test@example.com"},
		company: &models.Company{Name: "Test Co"},
	}
	db.invoice.ID = 1
	db.user.ID = 1

	svc := pay.NewService(db, &mockCrypto{})
	result, err := svc.Resolve(&pay.ResolveInput{InvoiceID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Invoice.Status != "paid" {
		t.Fatalf("expected status paid, got %s", result.Invoice.Status)
	}
	if result.EmailPayload == nil {
		t.Fatal("expected email payload")
	}
	if result.EmailPayload.AmountPaid != "$50.00" {
		t.Fatalf("expected $50.00, got %s", result.EmailPayload.AmountPaid)
	}
	if len(db.savedInvoices) != 1 {
		t.Fatalf("expected 1 save, got %d", len(db.savedInvoices))
	}
}

func TestResolve_ZeroAmount_NoEmail(t *testing.T) {
	db := &mockDB{
		invoice: &models.Invoice{Total: 0, UserID: 1, AppID: 1},
		user:    &models.User{Email: "test@example.com"},
		company: &models.Company{Name: "Test Co"},
	}
	db.invoice.ID = 1
	db.user.ID = 1

	svc := pay.NewService(db, &mockCrypto{})
	result, err := svc.Resolve(&pay.ResolveInput{InvoiceID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Invoice.Status != "paid" {
		t.Fatalf("expected status paid, got %s", result.Invoice.Status)
	}
	if result.EmailPayload != nil {
		t.Fatal("expected no email payload for zero amount")
	}
}

func TestResolve_ReactivatesPastDueSubscription(t *testing.T) {
	subID := uint(9)
	db := &mockDB{
		invoice: &models.Invoice{
			Total:          3000,
			UserID:         1,
			AppID:          1,
			InvoiceNumber:  "INV-009",
			SubscriptionID: &subID,
		},
		user:         &models.User{Email: "test@example.com"},
		company:      &models.Company{Name: "Test Co"},
		subscription: &models.Subscription{Status: "past_due"},
	}
	db.invoice.ID = 1
	db.user.ID = 1
	db.subscription.ID = subID

	svc := pay.NewService(db, &mockCrypto{})
	if _, err := svc.Resolve(&pay.ResolveInput{InvoiceID: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.subscription.Status != "active" {
		t.Fatalf("expected active, got %s", db.subscription.Status)
	}
	if len(db.savedSubscriptions) != 1 {
		t.Fatalf("expected 1 subscription save, got %d", len(db.savedSubscriptions))
	}
}
