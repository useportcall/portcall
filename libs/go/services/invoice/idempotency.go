package invoice

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func invoiceWindowStart(sub *models.Subscription) time.Time {
	now := time.Now()
	if sub.LastResetAt.IsZero() || sub.NextResetAt.After(now) {
		return sub.CreatedAt
	}
	return sub.NextResetAt
}
