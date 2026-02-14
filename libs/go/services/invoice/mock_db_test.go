package invoice_test

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

// mockDB implements dbx.IORM for invoice service tests.
type mockDB struct {
	dbx.IORM // embed to satisfy interface; only used methods are overridden

	subscription      *models.Subscription
	company           *models.Company
	invoiceCount      int64
	subItemIDs        []uint
	subscriptionItems map[uint]*models.SubscriptionItem
	billingMeters     []models.BillingMeter
	plans             map[uint]*models.Plan
	notFoundInvoice   bool // true = no existing invoice (idempotency check returns not found)

	// list support
	invoices      []models.Invoice
	subscriptions map[string]*models.Subscription // public_id -> subscription
	users         map[string]*models.User         // public_id -> user

	createdInvoices []any
	createdItems    []any
	savedInvoices   []any
}

func (m *mockDB) FindForID(id uint, dest any) error {
	switch d := dest.(type) {
	case *models.Subscription:
		if m.subscription == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.subscription
		return nil
	case *models.SubscriptionItem:
		if si, ok := m.subscriptionItems[id]; ok {
			*d = *si
			d.ID = id
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.Plan:
		if p, ok := m.plans[id]; ok {
			*d = *p
			d.ID = id
			return nil
		}
		return gorm.ErrRecordNotFound
	}
	return fmt.Errorf("mockDB.FindForID: unsupported type %T", dest)
}

func (m *mockDB) FindFirst(dest any, conds ...any) error {
	switch d := dest.(type) {
	case *models.Company:
		if m.company == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.company
		return nil
	case *models.Invoice:
		if m.notFoundInvoice {
			return gorm.ErrRecordNotFound
		}
		return nil // found = idempotent
	}
	return fmt.Errorf("mockDB.FindFirst: unsupported type %T", dest)
}

func (m *mockDB) Count(count *int64, dest any, query string, args ...any) error {
	*count = m.invoiceCount
	return nil
}

func (m *mockDB) ListIDs(table string, dest any, conds ...any) error {
	if ids, ok := dest.(*[]uint); ok {
		*ids = m.subItemIDs
		return nil
	}
	return nil
}

func (m *mockDB) Create(value any) error {
	m.createdInvoices = append(m.createdInvoices, value)
	return nil
}

func (m *mockDB) Save(dest any) error {
	m.savedInvoices = append(m.savedInvoices, dest)
	return nil
}

func (m *mockDB) Txn(fn func(orm dbx.IORM) error) error {
	return fn(m) // run in same mock — no real transaction
}

func (m *mockDB) GetForPublicID(appID uint, publicID string, dest any) error {
	switch d := dest.(type) {
	case *models.Subscription:
		if sub, ok := m.subscriptions[publicID]; ok {
			*d = *sub
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.User:
		if u, ok := m.users[publicID]; ok {
			*d = *u
			return nil
		}
		return gorm.ErrRecordNotFound
	}
	return fmt.Errorf("mockDB.GetForPublicID: unsupported type %T", dest)
}

func (m *mockDB) ListWithOrder(dest any, order string, conds ...any) error {
	if invoices, ok := dest.(*[]models.Invoice); ok {
		*invoices = m.invoices
		return nil
	}
	return fmt.Errorf("mockDB.ListWithOrder: unsupported type %T", dest)
}

func (m *mockDB) List(dest any, conds ...any) error {
	if meters, ok := dest.(*[]models.BillingMeter); ok {
		*meters = m.billingMeters
		return nil
	}
	return fmt.Errorf("mockDB.List: unsupported type %T", dest)
}
