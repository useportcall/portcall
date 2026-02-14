package payment_test

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/followup"
	pay "github.com/useportcall/portcall/libs/go/services/payment"
)

func TestProcessDunning_SetsFollowUpStage(t *testing.T) {
	subID := uint(12)
	db := &mockDB{
		invoice: &models.Invoice{
			Total:          3300,
			UserID:         1,
			AppID:          1,
			InvoiceNumber:  "INV-012",
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
	result, err := svc.ProcessDunning(&pay.DunningInput{InvoiceID: 1, Attempt: 1, MaxAttempts: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EmailPayload == nil || result.EmailPayload.FollowUpStage != followup.StageInvoiceFirst {
		t.Fatalf("expected first stage payload, got %+v", result.EmailPayload)
	}
}

func TestProcessDunning_NoRetryForcesFinal(t *testing.T) {
	subID := uint(13)
	db := &mockDB{
		invoice: &models.Invoice{
			Total:          3300,
			UserID:         1,
			AppID:          1,
			InvoiceNumber:  "INV-013",
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
	result, err := svc.ProcessDunning(&pay.DunningInput{InvoiceID: 1, Attempt: 1, MaxAttempts: 4, NoRetry: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.FinalAttempt || result.EmailPayload.Attempt != 4 {
		t.Fatalf("expected final attempt=4, got %+v", result.EmailPayload)
	}
}
