package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type CreateEntitlementsPayload struct {
	UserID uint `json:"user_id"`
	PlanID uint `json:"plan_id"`
}

func CreateEntitlements(c server.IContext) error {
	var p CreateEntitlementsPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var user models.User
	if err := c.DB().FindForID(p.UserID, &user); err != nil {
		return err
	}

	var planFeatures []models.PlanFeature
	if err := c.DB().List(&planFeatures, "plan_id = ?", p.PlanID); err != nil {
		return err
	}

	var entitlements []models.Entitlement
	if err := c.DB().List(&entitlements, "user_id = ?", p.UserID); err != nil {
		return err
	}

	lastReset := time.Now()

	// reset quota for all existing entitlements, TODO: look at more efficient way
	for _, ent := range entitlements {
		ent.Quota = 0
		ent.LastResetAt = &lastReset
		if err := c.DB().Save(&ent); err != nil {
			return err
		}
	}

	for _, pf := range planFeatures {
		var feature models.Feature
		if err := c.DB().FindForID(pf.FeatureID, &feature); err != nil {
			return err
		}

		entitlement := new(models.Entitlement)
		if err := c.DB().FindFirst(entitlement, "user_id = ? AND feature_public_id = ?", p.UserID, feature.PublicID); err != nil {
			if !dbx.IsRecordNotFoundError(err) {
				return err
			}

			nextReset, err := NextReset(lastReset, pf.Interval, lastReset)
			if err != nil {
				return err
			}

			log.Printf("🚀 Creating new entitlement for user %d, feature %d, plan %d", p.UserID, pf.FeatureID, p.PlanID)

			var planItem models.PlanItem
			if err := c.DB().FindForID(pf.PlanItemID, &planItem); err != nil {
				return err
			}

			var feature models.Feature
			if err := c.DB().FindForID(pf.FeatureID, &feature); err != nil {
				return err
			}

			entitlement.AppID = user.AppID
			entitlement.UserID = p.UserID
			entitlement.FeaturePublicID = feature.PublicID
			entitlement.Interval = pf.Interval
			entitlement.Quota = int64(pf.Quota)
			entitlement.Usage = 0
			entitlement.NextResetAt = nextReset
			entitlement.LastResetAt = &lastReset
			entitlement.AnchorAt = &lastReset
			entitlement.IsMetered = planItem.PricingModel != "fixed" // TODO: improve
			if err := c.DB().Create(&entitlement); err != nil {
				return err
			}
		} else {
			log.Println("🚀 Entitlement already exists, updating...")

			entitlement.LastResetAt = &lastReset
			entitlement.Interval = pf.Interval
			entitlement.Quota = int64(pf.Quota) // TODO: align these

			if err := c.DB().Save(entitlement); err != nil {
				return err
			}
		}
	}

	return nil
}
