package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type SubscriptionItem struct {
	ID        string    `db:"id" json:"-"`
	Quantity  int32     `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (si *SubscriptionItem) Set(subscriptionItem *models.SubscriptionItem) *SubscriptionItem {
	si.ID = subscriptionItem.PublicID
	si.Quantity = subscriptionItem.Quantity
	si.CreatedAt = subscriptionItem.CreatedAt
	si.UpdatedAt = subscriptionItem.UpdatedAt
	return si
}
