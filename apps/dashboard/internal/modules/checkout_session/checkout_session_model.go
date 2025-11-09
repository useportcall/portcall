package checkout_session

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/useportcall/portcall/apps/dashboard/internal/modules/address"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/company"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type CheckoutSession struct {
	ID                   string           `json:"id"`
	CreatedAt            time.Time        `json:"created_at"`
	ExpiresAt            time.Time        `json:"expires_at"`
	ExternalClientSecret *string          `json:"external_client_secret"`
	ExternalPublicKey    *string          `json:"external_public_key"`
	ExternalProvider     string           `json:"external_provider"`
	ExternalSessionID    string           `json:"external_session_id"`
	Plan                 *plan.Plan       `json:"plan"`
	RedirectURL          *string          `json:"redirect_url"`
	CancelURL            *string          `json:"cancel_url"`
	URL                  string           `json:"url"`
	BillingAddress       *address.Address `json:"billing_address"`
	CompanyAddress       *address.Address `json:"company_address"`
	Company              *company.Company `json:"company"`
}

func (cs *CheckoutSession) Set(checkoutSession *models.CheckoutSession) *CheckoutSession {
	cs.ID = checkoutSession.PublicID
	cs.CreatedAt = checkoutSession.CreatedAt
	cs.ExpiresAt = checkoutSession.ExpiresAt
	cs.ExternalClientSecret = &checkoutSession.ExternalClientSecret
	cs.ExternalPublicKey = &checkoutSession.ExternalPublicKey
	cs.ExternalProvider = checkoutSession.ExternalProvider
	cs.ExternalSessionID = checkoutSession.ExternalSessionID
	cs.RedirectURL = checkoutSession.RedirectURL
	cs.CancelURL = checkoutSession.CancelURL

	checkoutURL := os.Getenv("CHECKOUT_URL")
	if checkoutURL == "" {
		log.Println("CHECKOUT_URL environment variable is not set!")
	}

	cs.URL = fmt.Sprintf("%s?id=%s", checkoutURL, checkoutSession.PublicID)

	return cs
}

type CreateCheckoutSessionRequest struct {
	UserID      string `json:"user_id"`
	PlanID      string `json:"plan_id"`
	CancelURL   string `json:"cancel_url"`
	RedirectURL string `json:"redirect_url"`
}
