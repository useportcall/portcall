package payment_test

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

// mockDB implements dbx.IORM for payment service tests.
type mockDB struct {
	dbx.IORM // embed to satisfy interface

	invoice      *models.Invoice
	user         *models.User
	paymentMth   *models.PaymentMethod
	paymentMths  []models.PaymentMethod
	connection   *models.Connection
	company      *models.Company
	subscription *models.Subscription

	savedInvoices      []any
	savedSubscriptions []any
}

func (m *mockDB) FindForID(id uint, dest any) error {
	switch d := dest.(type) {
	case *models.Invoice:
		if m.invoice == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.invoice
		return nil
	case *models.User:
		if m.user == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.user
		return nil
	case *models.Subscription:
		if m.subscription == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.subscription
		return nil
	}
	return fmt.Errorf("mockDB.FindForID: unsupported type %T", dest)
}

func (m *mockDB) FindFirst(dest any, conds ...any) error {
	switch d := dest.(type) {
	case *models.PaymentMethod:
		if m.paymentMth == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.paymentMth
		return nil
	case *models.Connection:
		if m.connection == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.connection
		return nil
	case *models.Company:
		if m.company == nil {
			return gorm.ErrRecordNotFound
		}
		*d = *m.company
		return nil
	}
	return fmt.Errorf("mockDB.FindFirst: unsupported type %T", dest)
}

func (m *mockDB) Save(dest any) error {
	switch v := dest.(type) {
	case *models.Invoice:
		m.savedInvoices = append(m.savedInvoices, v)
		m.invoice = v
	case *models.Subscription:
		m.savedSubscriptions = append(m.savedSubscriptions, v)
		m.subscription = v
	}
	return nil
}

func (m *mockDB) Create(dest any) error {
	return nil
}

// mockCrypto implements cryptox.ICrypto for tests.
type mockCrypto struct{}

func (m *mockCrypto) Encrypt(data string) (string, error)            { return data, nil }
func (m *mockCrypto) Decrypt(data string) (string, error)            { return data, nil }
func (m *mockCrypto) CompareHash(hashed, plain string) (bool, error) { return true, nil }
