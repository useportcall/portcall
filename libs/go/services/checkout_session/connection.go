package checkout_session

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *service) resolveConnection(appID, defaultConnectionID uint) (*models.Connection, error) {
	if defaultConnectionID != 0 {
		var configured models.Connection
		err := s.db.FindForID(defaultConnectionID, &configured)
		if err == nil {
			return &configured, nil
		}
		if !dbx.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	var local models.Connection
	if err := s.db.FindFirst(&local, "app_id = ? AND source = ?", appID, "local"); err == nil {
		return &local, nil
	} else if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}

	var fallback models.Connection
	if err := s.db.FindFirst(&fallback, "app_id = ?", appID); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			return nil, NewValidationError("no payment connection configured")
		}
		return nil, err
	}
	return &fallback, nil
}
