package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Entitlement struct {
	ID          string     `json:"id"`
	Usage       int64      `json:"usage"`
	Quota       int64      `json:"quota"`
	Interval    string     `json:"interval"`
	Enabled     bool       `json:"enabled"`
	NextResetAt *time.Time `json:"next_reset_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (e *Entitlement) Set(entitlement *models.Entitlement) *Entitlement {
	e.ID = entitlement.FeaturePublicID
	e.Interval = entitlement.Interval
	e.Quota = entitlement.Quota
	e.Usage = entitlement.Usage
	e.NextResetAt = entitlement.NextResetAt
	e.CreatedAt = entitlement.CreatedAt
	e.UpdatedAt = entitlement.UpdatedAt

	if entitlement.Quota == -1 {
		e.Enabled = true
	} else {
		e.Enabled = entitlement.Usage < entitlement.Quota
	}

	return e
}
