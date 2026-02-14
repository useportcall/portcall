package payment

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

func findLatestPaymentMethod(db dbx.IORM, userID, appID uint) (*models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	if err := db.ListWithOrderAndLimit(
		&methods,
		"created_at DESC",
		1,
		"user_id = ? AND app_id = ?",
		userID,
		appID,
	); err != nil {
		return nil, err
	}
	if len(methods) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &methods[0], nil
}
