package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type IUpsertEntitlementsContext interface {
	GetStartUpsertEntitlementsPayload() (*StartUpsertEntitlementsPayload, error)
	ListPlanFeatureIDs(planID uint) ([]uint, error)
	EnqueueUpsertEntitlement(index int, values []uint, userID uint)
	GetUpsertEntitlementPayload() (*UpsertEntitlementPayload, error)
	FindPlanFeatureByID(id uint) (*models.PlanFeature, error)
	BuildEntitlement(userID uint, feature string, pf *models.PlanFeature) *models.Entitlement
	UpsertEntitlement(userID uint, feature string, create *models.Entitlement, update *models.Entitlement) error
}

type StartUpsertEntitlementsPayload struct {
	UserID uint `json:"user_id"`
	PlanID uint `json:"plan_id"`
}

func (c *Context) GetStartUpsertEntitlementsPayload() (*StartUpsertEntitlementsPayload, error) {
	var p StartUpsertEntitlementsPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Context) ListPlanFeatureIDs(planID uint) ([]uint, error) {
	var planFeatureIDs []uint
	if err := c.DB().ListIDs("plan_features", &planFeatureIDs, "plan_id = ?", planID); err != nil {
		return nil, err
	}

	return planFeatureIDs, nil
}

type UpsertEntitlementPayload struct {
	Index  int    `json:"index"`
	Values []uint `json:"values"`
	UserID uint   `json:"user_id"`
}

func (c *Context) EnqueueUpsertEntitlement(index int, values []uint, userID uint) {
	payload := UpsertEntitlementPayload{
		Index:  index,
		Values: values,
		UserID: userID,
	}

	if err := c.Queue().Enqueue("upsert_entitlement", payload, "billing_queue"); err != nil {
		log.Printf("Error enqueueing upsert_entitlement task: %v", err)
	}

	log.Printf("Enqueued upsert_entitlement task for user ID %d, index %d", userID, index)
}

func (c *Context) GetUpsertEntitlementPayload() (*UpsertEntitlementPayload, error) {
	var p UpsertEntitlementPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Context) FindPlanFeatureByID(id uint) (*models.PlanFeature, error) {
	var planFeature models.PlanFeature
	if err := c.DB().FindForID(id, &planFeature); err != nil {
		return nil, err
	}
	return &planFeature, nil
}

func (c *Context) BuildEntitlement(userID uint, feature string, pf *models.PlanFeature) *models.Entitlement {
	lastReset := time.Now()

	entitlement := new(models.Entitlement)
	entitlement.AppID = pf.AppID
	entitlement.UserID = userID
	entitlement.FeaturePublicID = feature
	entitlement.Interval = pf.Interval
	entitlement.Quota = int64(pf.Quota)
	entitlement.Usage = 0
	entitlement.LastResetAt = &lastReset
	entitlement.AnchorAt = &lastReset
	// entitlement.IsMetered = feature.IsMetered // TODO: look into

	return entitlement
}

func (c *Context) UpsertEntitlement(userID uint, feature string, create *models.Entitlement, update *models.Entitlement) error {
	entitlement := new(models.Entitlement)

	err := c.DB().FindFirst(entitlement, "user_id = ? AND feature_public_id = ?", userID, feature)

	if !dbx.IsRecordNotFoundError(err) {
		return err
	}

	if err == nil {
		entitlement.Quota = update.Quota
		entitlement.Interval = update.Interval
		return c.DB().Save(entitlement)
	}

	// Create new entitlement
	if err := c.DB().Create(create); err != nil {
		return err
	}

	return nil
}
