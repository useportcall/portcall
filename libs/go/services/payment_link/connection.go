package payment_link

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *service) ensurePaymentConnection(appID uint) error {
	var config models.AppConfig
	if err := s.db.FindFirst(&config, "app_id = ?", appID); err == nil {
		if config.DefaultConnectionID != 0 {
			var configured models.Connection
			if err := s.db.FindForID(config.DefaultConnectionID, &configured); err == nil {
				return nil
			}
		}
	} else if !dbx.IsRecordNotFoundError(err) {
		return err
	}

	var fallback models.Connection
	if err := s.db.FindFirst(&fallback, "app_id = ?", appID); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			return NewValidationError("no payment connection configured")
		}
		return err
	}
	return nil
}
