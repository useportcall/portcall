package payment_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/services/payment"
	"gorm.io/gorm"
)

func TestCreatePaymentMethod_Success(t *testing.T) {
	db := &mockDB{
		// FindFirst returns not found so the method is created
	}
	svc := payment.NewService(db, &mockCrypto{})

	result, err := svc.CreatePaymentMethod(&payment.CreateMethodInput{
		AppID:                   1,
		UserID:                  2,
		PlanID:                  10,
		ExternalPaymentMethodID: "pm_stripe_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AppID != 1 || result.UserID != 2 || result.PlanID != 10 {
		t.Fatal("result fields mismatch")
	}
	if result.PaymentMethod == nil {
		t.Fatal("expected payment method in result")
	}
	if result.PaymentMethod.ExternalID != "pm_stripe_123" {
		t.Fatalf("expected pm_stripe_123, got %s", result.PaymentMethod.ExternalID)
	}
}

func TestCreatePaymentMethod_AlreadyExists(t *testing.T) {
	db := &mockDB{
		paymentMth: &models.PaymentMethod{
			Model:      gorm.Model{ID: 99},
			UserID:     2,
			ExternalID: "pm_stripe_123",
		},
	}
	svc := payment.NewService(db, &mockCrypto{})

	result, err := svc.CreatePaymentMethod(&payment.CreateMethodInput{
		AppID:                   1,
		UserID:                  2,
		PlanID:                  10,
		ExternalPaymentMethodID: "pm_stripe_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AppID != 1 || result.UserID != 2 || result.PlanID != 10 {
		t.Fatal("result fields mismatch")
	}
}
