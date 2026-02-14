package entitlement

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

type upsertMockDB struct {
	dbx.IORM
	saved models.Entitlement
}

func (m *upsertMockDB) FindForID(_ uint, dest any) error {
	switch d := dest.(type) {
	case *models.PlanFeature:
		*d = models.PlanFeature{AppID: 1, FeatureID: 10, Interval: "month", Quota: 50}
	case *models.Feature:
		*d = models.Feature{PublicID: "feat_tokens", IsMetered: true}
	}
	return nil
}

func (m *upsertMockDB) FindFirst(dest any, _ ...any) error {
	if ent, ok := dest.(*models.Entitlement); ok {
		*ent = models.Entitlement{
			Model: gorm.Model{ID: 99}, UserID: 1, FeaturePublicID: "feat_tokens",
			Usage: 27, Quota: 10, Interval: "week", IsMetered: false,
		}
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *upsertMockDB) Save(dest any) error {
	m.saved = *(dest.(*models.Entitlement))
	return nil
}

func TestUpsert_ResetsExistingUsageAndMeteredState(t *testing.T) {
	db := &upsertMockDB{}
	_, err := NewService(db).Upsert(&UpsertInput{UserID: 1, Index: 0, Values: []uint{1}})
	if err != nil {
		t.Fatal(err)
	}
	if db.saved.Usage != 0 || db.saved.Quota != 50 {
		t.Fatalf("expected usage reset and new quota, got usage=%d quota=%d", db.saved.Usage, db.saved.Quota)
	}
	if !db.saved.IsMetered {
		t.Fatal("expected is_metered=true after upsert")
	}
}
