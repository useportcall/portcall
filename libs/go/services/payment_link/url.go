package payment_link

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func buildPaymentLinkURL(baseURL string, link *models.PaymentLink, crypto cryptox.ICrypto) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" || link == nil {
		return "", nil
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid CHECKOUT_URL: %w", err)
	}
	query := parsed.Query()
	query.Set("pl", link.PublicID)
	if crypto != nil {
		token, err := cryptox.CreatePaymentLinkToken(crypto, link.PublicID, link.ExpiresAt)
		if err != nil {
			return "", err
		}
		query.Set("pt", token)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
