package subscription_test

import (
	"fmt"
	"strings"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

type mockDB struct {
	dbx.IORM
	subscriptions   map[uint]*models.Subscription
	subItems        map[uint]*models.SubscriptionItem
	plans           map[uint]*models.Plan
	planItems       map[uint]*models.PlanItem
	planItemsByPlan map[uint]*models.PlanItem
	users           map[uint]*models.User
	freePlans       []models.Plan
	planItemIDs     []uint
	created         []any
	saved           []any
	deleted         []any
}

func newMockDB() *mockDB {
	return &mockDB{
		subscriptions:   make(map[uint]*models.Subscription),
		subItems:        make(map[uint]*models.SubscriptionItem),
		plans:           make(map[uint]*models.Plan),
		planItems:       make(map[uint]*models.PlanItem),
		planItemsByPlan: make(map[uint]*models.PlanItem),
		users:           make(map[uint]*models.User),
	}
}

func (m *mockDB) FindForID(id uint, dest any) error {
	switch d := dest.(type) {
	case *models.Subscription:
		if s, ok := m.subscriptions[id]; ok {
			*d = *s
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.Plan:
		if p, ok := m.plans[id]; ok {
			*d = *p
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.User:
		if u, ok := m.users[id]; ok {
			*d = *u
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.PlanItem:
		if pi, ok := m.planItems[id]; ok {
			*d = *pi
			return nil
		}
		return gorm.ErrRecordNotFound
	}
	return fmt.Errorf("mockDB.FindForID: unsupported %T", dest)
}

func (m *mockDB) FindFirst(dest any, conds ...any) error {
	switch d := dest.(type) {
	case *models.Subscription:
		for _, s := range m.subscriptions {
			*d = *s
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.PlanItem:
		if len(conds) >= 2 {
			if planID, ok := conds[1].(uint); ok {
				if pi, found := m.planItemsByPlan[planID]; found {
					*d = *pi
					return nil
				}
			}
		}
		for _, pi := range m.planItems {
			*d = *pi
			return nil
		}
		return gorm.ErrRecordNotFound
	}
	return fmt.Errorf("mockDB.FindFirst: unsupported %T", dest)
}

func (m *mockDB) List(dest any, conds ...any) error {
	if plans, ok := dest.(*[]models.Plan); ok {
		*plans = m.freePlans
		return nil
	}
	return nil
}

func (m *mockDB) ListIDs(table string, dest any, conds ...any) error {
	if ids, ok := dest.(*[]uint); ok {
		if table == "plans" {
			query, _ := conds[0].(string)
			if strings.Contains(query, "is_free") && len(m.freePlans) > 0 {
				for _, plan := range m.freePlans {
					*ids = append(*ids, plan.ID)
				}
				return nil
			}
			for id, plan := range m.plans {
				if strings.Contains(query, "is_free") && !plan.IsFree {
					continue
				}
				*ids = append(*ids, id)
			}
			return nil
		}
		*ids = m.planItemIDs
		return nil
	}
	return nil
}

func (m *mockDB) Create(value any) error {
	switch v := value.(type) {
	case *models.Subscription:
		if v.ID == 0 {
			v.ID = uint(len(m.subscriptions) + 1)
		}
		m.subscriptions[v.ID] = v
	case *models.SubscriptionItem:
		if v.ID == 0 {
			v.ID = uint(len(m.subItems) + 1)
		}
		m.subItems[v.ID] = v
	}
	m.created = append(m.created, value)
	return nil
}

func (m *mockDB) Save(dest any) error {
	m.saved = append(m.saved, dest)
	return nil
}

func (m *mockDB) Delete(_ any, _ any, _ ...any) error {
	m.deleted = append(m.deleted, true)
	return nil
}

func (m *mockDB) Txn(fn func(dbx.IORM) error) error {
	subscriptions := cloneSubscriptions(m.subscriptions)
	subItems := cloneSubItems(m.subItems)
	createdLen := len(m.created)
	savedLen := len(m.saved)
	deletedLen := len(m.deleted)
	if err := fn(m); err != nil {
		m.subscriptions = subscriptions
		m.subItems = subItems
		m.created = m.created[:createdLen]
		m.saved = m.saved[:savedLen]
		m.deleted = m.deleted[:deletedLen]
		return err
	}
	return nil
}

func cloneSubscriptions(
	src map[uint]*models.Subscription,
) map[uint]*models.Subscription {
	dst := make(map[uint]*models.Subscription, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func cloneSubItems(
	src map[uint]*models.SubscriptionItem,
) map[uint]*models.SubscriptionItem {
	dst := make(map[uint]*models.SubscriptionItem, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
