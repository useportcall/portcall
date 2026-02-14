package invoice

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/routerx"
	inv "github.com/useportcall/portcall/libs/go/services/invoice"
)

// ListInvoices handles GET /v1/invoices.
func ListInvoices(c *routerx.Context) {
	svc := inv.NewService(c.DB())
	result, err := svc.List(&inv.ListInput{
		AppID:          c.AppID(),
		SubscriptionID: c.Query("subscription_id"),
		UserID:         c.Query("user_id"),
	})
	if err != nil {
		c.ServerError("Failed to list invoices", err)
		return
	}

	response := make([]apix.Invoice, len(result.Invoices))
	for i, invoice := range result.Invoices {
		response[i].Set(&invoice)
	}
	c.OK(response)
}
