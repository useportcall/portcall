package apix

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type CheckoutSession struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	ExpiresAt            time.Time `json:"expires_at"`
	ExternalClientSecret *string   `json:"external_client_secret"`
	ExternalPublicKey    *string   `json:"external_public_key"`
	ExternalProvider     string    `json:"external_provider"`
	ExternalSessionID    string    `json:"external_session_id"`
	Plan                 *Plan     `json:"plan"`
	RedirectURL          *string   `json:"redirect_url"`
	CancelURL            *string   `json:"cancel_url"`
	URL                  string    `json:"url"`
	BillingAddress       *Address  `json:"billing_address"`
	CompanyAddress       *Address  `json:"company_address"`
	Company              *Company  `json:"company"`
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
