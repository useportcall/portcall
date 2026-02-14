package paymentx

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
)

type localPaymentClient struct{}

func (c *localPaymentClient) CreateCustomer(email, name string) (string, error) {
	return dbx.GenPublicID("cust"), nil
}

func (c *localPaymentClient) CreateCheckoutSession(customerID string) (string, string, error) {
	return dbx.GenPublicID("rand"), dbx.GenPublicID("cls"), nil
}

func (c *localPaymentClient) CreateCharge(customerID string, total int64, currency string, paymentMethodID string, metadata map[string]string) error {
	log.Printf("Local payment charge created for customer (%s) with amount: %d %s", customerID, total, currency)
	return nil
}

func (c *localPaymentClient) Verify() error { return nil }

func (c *localPaymentClient) RegisterWebhook(webhookURL string) (*WebhookRegistration, error) {
	log.Printf("Local payment webhook registered: %s", webhookURL)
	return &WebhookRegistration{
		EndpointID: dbx.GenPublicID("we"),
		Secret:     dbx.GenPublicID("whsec"),
	}, nil
}

func (c *localPaymentClient) DeleteWebhook(endpointID string) error {
	log.Printf("Local payment webhook deleted: %s", endpointID)
	return nil
}
