package entitlement

import (
	"log"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// ResetAll resets all entitlements for a user in a single transaction.
// It sets usage to 0 and updates LastResetAt for each entitlement.
func (s *service) ResetAll(input *ResetAllInput) (*ResetAllResult, error) {
	log.Printf("Processing ResetAll for user ID %d", input.UserID)

	entitlements, err := findEntitlements(s.db, input.UserID)
	if err != nil {
		return nil, err
	}

	if len(entitlements) == 0 {
		log.Printf("No entitlements found for user %d", input.UserID)
		return &ResetAllResult{ResetCount: 0}, nil
	}

	if err := resetEntitlements(s.db, entitlements); err != nil {
		return nil, err
	}

	log.Printf("Reset %d entitlements for user %d", len(entitlements), input.UserID)
	return &ResetAllResult{ResetCount: len(entitlements)}, nil
}

func findEntitlements(db dbx.IORM, userID uint) ([]models.Entitlement, error) {
	var entitlements []models.Entitlement
	if err := db.List(&entitlements, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return entitlements, nil
}

func resetEntitlements(db dbx.IORM, entitlements []models.Entitlement) error {
	return db.Txn(func(tx dbx.IORM) error {
		now := time.Now()
		for i := range entitlements {
			entitlements[i].Usage = 0
			entitlements[i].LastResetAt = &now
			entitlements[i].NextResetAt = nextResetAt(entitlements[i].Interval, now)
			if err := tx.Save(&entitlements[i]); err != nil {
				return err
			}
		}
		return nil
	})
}
