package payment_test

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	pay "github.com/useportcall/portcall/libs/go/services/payment"
)

func TestProcessDunning_IntermediateAttempt(t *testing.T) {
	subID := uint(7)
	db := &mockDB{
		invoice: &models.Invoice{
			Total:          2500,
			UserID:         1,
			AppID:          1,
			InvoiceNumber:  "INV-007",
			DueBy:          time.Now().AddDate(0, 0, 7),
			Status:         "issued",
			SubscriptionID: &subID,
		},
		user:         &models.User{Email: "test@example.com"},
		company:      &models.Company{Name: "Test Co"},
		subscription: &models.Subscription{Status: "active"},
	}
	db.invoice.ID = 1
	db.subscription.ID = subID

	svc := pay.NewService(db, &mockCrypto{})
	result, err := svc.ProcessDunning(&pay.DunningInput{
		InvoiceID: 1, Attempt: 1, MaxAttempts: 3, FailureReason: "card declined",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.FinalAttempt {
		t.Fatal("expected non-final attempt")
	}
	if db.invoice.Status != "past_due" {
		t.Fatalf("expected past_due, got %s", db.invoice.Status)
	}
	if db.subscription.Status != "past_due" {
		t.Fatalf("expected past_due subscription, got %s", db.subscription.Status)
	}
}

func TestProcessDunning_FinalAttempt(t *testing.T) {
	subID := uint(8)
	db := &mockDB{
		invoice: &models.Invoice{
			Total:          1200,
			UserID:         1,
			AppID:          1,
			InvoiceNumber:  "INV-008",
			DueBy:          time.Now().AddDate(0, 0, 7),
			Status:         "past_due",
			SubscriptionID: &subID,
		},
		user:         &models.User{Email: "test@example.com"},
		company:      &models.Company{Name: "Test Co"},
		subscription: &models.Subscription{Status: "past_due"},
	}
	db.invoice.ID = 1
	db.subscription.ID = subID

	svc := pay.NewService(db, &mockCrypto{})
	result, err := svc.ProcessDunning(&pay.DunningInput{
		InvoiceID: 1, Attempt: 3, MaxAttempts: 3, FailureReason: "insufficient funds",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.FinalAttempt {
		t.Fatal("expected final attempt")
	}
	if db.invoice.Status != "uncollectible" {
		t.Fatalf("expected uncollectible, got %s", db.invoice.Status)
	}
	if db.subscription.Status != "canceled" {
		t.Fatalf("expected canceled subscription, got %s", db.subscription.Status)
	}
}

func TestProcessDunning_DuplicateFirstAttempt_SkipsEmail(t *testing.T) {
	subID := uint(11)
	db := &mockDB{
		invoice: &models.Invoice{
			Total:          1200,
			UserID:         1,
			AppID:          1,
			InvoiceNumber:  "INV-011",
			DueBy:          time.Now().AddDate(0, 0, 7),
			Status:         "past_due",
			SubscriptionID: &subID,
		},
		user:         &models.User{Email: "test@example.com"},
		company:      &models.Company{Name: "Test Co"},
		subscription: &models.Subscription{Status: "past_due"},
	}
	db.invoice.ID = 1
	db.subscription.ID = subID

	svc := pay.NewService(db, &mockCrypto{})
	result, err := svc.ProcessDunning(&pay.DunningInput{
		InvoiceID: 1, Attempt: 1, MaxAttempts: 4, FailureReason: "duplicate event",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EmailPayload != nil {
		t.Fatal("expected no email for duplicate first-attempt webhook event")
	}
}
