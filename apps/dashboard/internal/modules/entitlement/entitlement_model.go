package entitlement

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Entitlement struct {
	ID          string     `json:"id"`
	Usage       int64      `json:"usage"`
	Quota       int64      `json:"quota"`
	Interval    string     `json:"interval"`
	NextResetAt *time.Time `json:"next_reset_at"`
	Feature     any        `json:"feature"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (e *Entitlement) Set(entitlement *models.Entitlement) *Entitlement {
	e.ID = entitlement.PublicID
	e.Interval = entitlement.Interval
	e.Quota = entitlement.Quota
	e.Usage = entitlement.Usage
	e.NextResetAt = entitlement.NextResetAt
	e.CreatedAt = entitlement.CreatedAt
	e.UpdatedAt = entitlement.UpdatedAt
	return e
}
