package subscription

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

type rollbackMockDB struct {
	dbx.IORM
	plans map[uint]models.Plan
	items map[uint]models.PlanItem
}

func (m *rollbackMockDB) ListIDs(table string, dest any, conds ...any) error {
	ids := dest.(*[]uint)
	if table != "plans" {
		return nil
	}
	query, _ := conds[0].(string)
	for id, plan := range m.plans {
		if query == "app_id = ? AND is_free = ?" && !plan.IsFree {
			continue
		}
		*ids = append(*ids, id)
	}
	return nil
}

func (m *rollbackMockDB) FindForID(id uint, dest any) error {
	if p, ok := dest.(*models.Plan); ok {
		*p = m.plans[id]
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *rollbackMockDB) FindFirst(dest any, conds ...any) error {
	if item, ok := dest.(*models.PlanItem); ok {
		planID, _ := conds[1].(uint)
		if pi, found := m.items[planID]; found {
			*item = pi
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func TestFindRollbackPlanID_PrefersFreePlan(t *testing.T) {
	db := &rollbackMockDB{
		plans: map[uint]models.Plan{
			10: {Model: gorm.Model{ID: 10}, IsFree: false},
			20: {Model: gorm.Model{ID: 20}, IsFree: true},
		},
	}
	cur := uint(10)
	id, err := findRollbackPlanID(db, 1, &cur)
	if err != nil || id == nil || *id != 20 {
		t.Fatalf("expected free rollback plan 20, got %v err=%v", id, err)
	}
}

func TestFindRollbackPlanID_UsesCheapestFixedPlan(t *testing.T) {
	db := &rollbackMockDB{
		plans: map[uint]models.Plan{
			10: {Model: gorm.Model{ID: 10}, IsFree: false},
			30: {Model: gorm.Model{ID: 30}, IsFree: false},
			40: {Model: gorm.Model{ID: 40}, IsFree: false},
		},
		items: map[uint]models.PlanItem{
			30: {PlanID: 30, PricingModel: "fixed", UnitAmount: 500},
			40: {PlanID: 40, PricingModel: "fixed", UnitAmount: 100},
		},
	}
	cur := uint(10)
	id, err := findRollbackPlanID(db, 1, &cur)
	if err != nil || id == nil || *id != 40 {
		t.Fatalf("expected cheapest rollback plan 40, got %v err=%v", id, err)
	}
}
