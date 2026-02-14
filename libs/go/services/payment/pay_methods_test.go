package payment

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

type latestPMDB struct {
	dbx.IORM
	methods []models.PaymentMethod
}

func (m latestPMDB) ListWithOrderAndLimit(dest any, order string, limit int, conds ...any) error {
	target, ok := dest.(*[]models.PaymentMethod)
	if !ok {
		return gorm.ErrInvalidData
	}
	n := len(m.methods)
	if limit > 0 && limit < n {
		n = limit
	}
	*target = append([]models.PaymentMethod(nil), m.methods[:n]...)
	return nil
}

func TestFindLatestPaymentMethod_ReturnsFirstRecord(t *testing.T) {
	db := latestPMDB{
		methods: []models.PaymentMethod{
			{ExternalID: "pm_new"},
			{ExternalID: "pm_old"},
		},
	}

	pm, err := findLatestPaymentMethod(db, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pm.ExternalID != "pm_new" {
		t.Fatalf("expected pm_new, got %s", pm.ExternalID)
	}
}

func TestFindLatestPaymentMethod_EmptyReturnsNotFound(t *testing.T) {
	_, err := findLatestPaymentMethod(latestPMDB{}, 1, 1)
	if err != gorm.ErrRecordNotFound {
		t.Fatalf("expected record not found, got %v", err)
	}
}
