package payment_test

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	pay "github.com/useportcall/portcall/libs/go/services/payment"
)

func TestBraintreeDeclinedTriggersDunning(t *testing.T) {
	svc := pay.NewService(&mockDB{}, &mockCrypto{})
	result, err := svc.ProcessBraintreeWebhook(
		&pay.BraintreeWebhookInput{
			Kind:          "transaction_settlement_declined",
			OrderID:       "portcall_invoice_id=99",
			FailureCount:  3,
			FailureReason: "processor_declined",
		},
	)
	if err != nil {
		t.Fatalf("ProcessBraintreeWebhook error: %v", err)
	}
	if result.Action != "process_stripe_payment_failure" {
		t.Fatalf("action = %s, want dunning", result.Action)
	}
	if result.Failure == nil {
		t.Fatal("expected failure payload")
	}
	if result.Failure.InvoiceID != 99 {
		t.Fatalf("invoice_id = %d, want 99", result.Failure.InvoiceID)
	}

	subID := uint(50)
	db := &mockDB{
		invoice: &models.Invoice{
			Total: 5000, UserID: 1, AppID: 1,
			InvoiceNumber: "INV-099", Status: "past_due",
			DueBy:          time.Now().AddDate(0, 0, -1),
			SubscriptionID: &subID,
		},
		user:         &models.User{Email: "u@test.com"},
		company:      &models.Company{Name: "Co"},
		subscription: &models.Subscription{Status: "past_due"},
	}
	db.invoice.ID = 99
	db.subscription.ID = subID

	svc2 := pay.NewService(db, &mockCrypto{})
	dunning, err := svc2.ProcessDunning(&pay.DunningInput{
		InvoiceID:     result.Failure.InvoiceID,
		Attempt:       3,
		MaxAttempts:   3,
		FailureReason: result.Failure.FailureReason,
	})
	if err != nil {
		t.Fatalf("ProcessDunning error: %v", err)
	}
	if !dunning.FinalAttempt {
		t.Fatal("expected final attempt")
	}
	if db.invoice.Status != "uncollectible" {
		t.Fatalf("invoice = %s, want uncollectible", db.invoice.Status)
	}
	if db.subscription.Status != "canceled" {
		t.Fatalf("sub = %s, want canceled", db.subscription.Status)
	}
}
