package entitlement_test

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

// mockDB implements dbx.IORM for entitlement service tests.
type mockDB struct {
	dbx.IORM // embed to satisfy interface

	entitlements []models.Entitlement
	planFeatures map[uint]*models.PlanFeature
	features     map[uint]*models.Feature

	savedEntitlements   []any
	createdEntitlements []any
	txnCalled           bool
}

func (m *mockDB) List(dest any, conds ...any) error {
	if ents, ok := dest.(*[]models.Entitlement); ok {
		*ents = m.entitlements
		return nil
	}
	return fmt.Errorf("mockDB.List: unsupported type %T", dest)
}

func (m *mockDB) FindForID(id uint, dest any) error {
	switch d := dest.(type) {
	case *models.PlanFeature:
		if pf, ok := m.planFeatures[id]; ok {
			*d = *pf
			d.ID = id
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.Feature:
		if f, ok := m.features[id]; ok {
			*d = *f
			d.ID = id
			return nil
		}
		return gorm.ErrRecordNotFound
	}
	return fmt.Errorf("mockDB.FindForID: unsupported type %T", dest)
}

func (m *mockDB) FindFirst(dest any, conds ...any) error {
	if _, ok := dest.(*models.Entitlement); ok {
		return gorm.ErrRecordNotFound
	}
	return fmt.Errorf("mockDB.FindFirst: unsupported type %T", dest)
}

func (m *mockDB) ListIDs(table string, dest any, conds ...any) error {
	return nil
}

func (m *mockDB) Create(value any) error {
	m.createdEntitlements = append(m.createdEntitlements, value)
	return nil
}

func (m *mockDB) Save(dest any) error {
	m.savedEntitlements = append(m.savedEntitlements, dest)
	return nil
}

func (m *mockDB) Txn(fn func(orm dbx.IORM) error) error {
	m.txnCalled = true
	return fn(m)
}
