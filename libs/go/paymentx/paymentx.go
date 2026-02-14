package paymentx

import (
	"fmt"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// WebhookRegistration holds the result of registering a webhook endpoint
type WebhookRegistration struct {
	EndpointID string // The external provider's webhook endpoint ID (for management/deletion)
	Secret     string // The webhook signing secret
}

type IPaymentClient interface {
	CreateCustomer(email, name string) (string, error)
	CreateCheckoutSession(customerID string) (string, string, error)
	CreateCharge(customerID string, total int64, currency string, paymentMethodID string, metadata map[string]string) error
	Verify() error
	RegisterWebhook(webhookURL string) (*WebhookRegistration, error)
	DeleteWebhook(endpointID string) error
}

func New(connection *models.Connection, crypto cryptox.ICrypto) (IPaymentClient, error) {
	switch connection.Source {
	case "local":
		return &localPaymentClient{}, nil
	case "stripe":
		secretKey, err := crypto.Decrypt(connection.EncryptedKey)
		if err != nil {
			return nil, err
		}
		return newStripePaymentClient(secretKey), nil
	case "braintree":
		return newBraintreePaymentClient(connection, crypto)
	default:
		return nil, fmt.Errorf("unsupported payment source: %s", connection.Source)
	}
}
