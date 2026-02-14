package entitlement

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

type startMockDB struct {
	dbx.IORM
	deleted []uint
}

func (m *startMockDB) ListIDs(table string, dest any, _ ...any) error {
	ids := dest.(*[]uint)
	switch table {
	case "plan_features":
		*ids = []uint{1}
	case "entitlements":
		*ids = []uint{100, 101}
	}
	return nil
}

func (m *startMockDB) FindForID(id uint, dest any) error {
	switch d := dest.(type) {
	case *models.PlanFeature:
		*d = models.PlanFeature{FeatureID: 10}
	case *models.Feature:
		*d = models.Feature{PublicID: "allowed"}
	case *models.Entitlement:
		if id == 100 {
			*d = models.Entitlement{Model: gorm.Model{ID: id}, FeaturePublicID: "allowed"}
		} else {
			*d = models.Entitlement{Model: gorm.Model{ID: id}, FeaturePublicID: "stale"}
		}
	}
	return nil
}

func (m *startMockDB) Delete(_ any, _ any, args ...any) error {
	id, _ := args[0].(uint)
	m.deleted = append(m.deleted, id)
	return nil
}

func TestStartUpsert_RemovesStaleEntitlements(t *testing.T) {
	db := &startMockDB{}
	result, err := NewService(db).StartUpsert(&StartUpsertInput{UserID: 1, PlanID: 2})
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasMore || len(result.PlanFeatureIDs) != 1 {
		t.Fatalf("unexpected start result: %+v", result)
	}
	if len(db.deleted) != 1 || db.deleted[0] != 101 {
		t.Fatalf("expected stale entitlement 101 to be deleted, got %v", db.deleted)
	}
}
