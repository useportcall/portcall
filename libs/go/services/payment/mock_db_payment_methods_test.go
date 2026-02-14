package payment_test

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (m *mockDB) ListWithOrderAndLimit(dest any, order string, limit int, conds ...any) error {
	switch d := dest.(type) {
	case *[]models.PaymentMethod:
		if len(m.paymentMths) > 0 {
			n := len(m.paymentMths)
			if limit > 0 && limit < n {
				n = limit
			}
			*d = append([]models.PaymentMethod(nil), m.paymentMths[:n]...)
			return nil
		}
		if m.paymentMth != nil {
			*d = []models.PaymentMethod{*m.paymentMth}
			return nil
		}
		*d = []models.PaymentMethod{}
		return nil
	}
	return fmt.Errorf("mockDB.ListWithOrderAndLimit: unsupported type %T", dest)
}
