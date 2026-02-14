package payment_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	pay "github.com/useportcall/portcall/libs/go/services/payment"
)

func TestPay_ZeroAmount_SkipsPayment(t *testing.T) {
	db := &mockDB{
		invoice: &models.Invoice{Total: 0, UserID: 1, AppID: 1},
	}
	db.invoice.ID = 1

	svc := pay.NewService(db, &mockCrypto{})
	result, err := svc.Pay(&pay.PayInput{InvoiceID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Invoice.Total != 0 {
		t.Fatal("expected zero total")
	}
}

func TestPay_NegativeAmount_SkipsPayment(t *testing.T) {
	db := &mockDB{
		invoice: &models.Invoice{Total: -100, UserID: 1, AppID: 1},
	}
	db.invoice.ID = 1

	svc := pay.NewService(db, &mockCrypto{})
	result, err := svc.Pay(&pay.PayInput{InvoiceID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Invoice.Total != -100 {
		t.Fatalf("expected -100, got %d", result.Invoice.Total)
	}
}
