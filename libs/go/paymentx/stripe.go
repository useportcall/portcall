package paymentx

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type stripePaymentClient struct {
	stp *client.API
}

func newStripePaymentClient(secretKey string) IPaymentClient {
	stp := &client.API{}
	stp.Init(secretKey, nil)
	return stripePaymentClient{stp: stp}
}

func (c stripePaymentClient) CreateCustomer(email, name string) (string, error) {
	params := &stripe.CustomerParams{Email: stripe.String(email), Name: stripe.String(name)}
	customer, err := c.stp.Customers.New(params)
	if err != nil {
		return "", fmt.Errorf("error creating Stripe customer: %w", err)
	}
	return customer.ID, nil
}

func (c stripePaymentClient) CreateCheckoutSession(customerID string) (string, string, error) {
	params := &stripe.SetupIntentParams{
		Customer:           stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Usage:              stripe.String(string(stripe.SetupIntentUsageOffSession)),
	}
	setupIntent, err := c.stp.SetupIntents.New(params)
	if err != nil {
		return "", "", err
	}
	return setupIntent.ID, setupIntent.ClientSecret, nil
}

func (c stripePaymentClient) CreateCharge(customerID string, total int64, currency string, paymentMethodID string, metadata map[string]string) error {
	params := &stripe.PaymentIntentParams{
		Params:             stripe.Params{Metadata: metadata},
		Customer:           stripe.String(customerID),
		Amount:             stripe.Int64(total),
		Currency:           stripe.String(currency),
		PaymentMethod:      stripe.String(paymentMethodID),
		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodManual)),
		Confirm:            stripe.Bool(true),
		OffSession:         stripe.Bool(true),
	}
	if invoiceID := metadata["portcall_invoice_id"]; invoiceID != "" {
		params.SetIdempotencyKey("portcall-invoice-" + invoiceID)
	}
	paymentIntent, err := c.stp.PaymentIntents.New(params)
	if err != nil {
		return fmt.Errorf("failed to create payment intent for customer (%s): %w", customerID, err)
	}
	if paymentIntent.Status == stripe.PaymentIntentStatusRequiresAction {
		return fmt.Errorf("payment intent for customer (%s) requires action", customerID)
	}
	if paymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
		return fmt.Errorf("payment intent for customer (%s) status is not succeeded: %s", customerID, paymentIntent.Status)
	}
	log.Printf("Payment intent for customer (%s) created successfully with ID: %s", customerID, paymentIntent.ID)
	return nil
}

func (c stripePaymentClient) Verify() error {
	_, err := c.stp.Account.Get()
	if err != nil {
		return fmt.Errorf("failed to verify Stripe payment client: %w", err)
	}
	return nil
}

func (c stripePaymentClient) RegisterWebhook(webhookURL string) (*WebhookRegistration, error) {
	params := &stripe.WebhookEndpointParams{
		URL:           stripe.String(webhookURL),
		EnabledEvents: stripe.StringSlice(stripeWebhookEvents),
	}
	endpoint, err := c.stp.WebhookEndpoints.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to register Stripe webhook endpoint: %w", err)
	}
	log.Printf("Stripe webhook endpoint registered: ID=%s, URL=%s", endpoint.ID, webhookURL)
	return &WebhookRegistration{EndpointID: endpoint.ID, Secret: endpoint.Secret}, nil
}

func (c stripePaymentClient) DeleteWebhook(endpointID string) error {
	_, err := c.stp.WebhookEndpoints.Del(endpointID, nil)
	if err != nil {
		return fmt.Errorf("failed to delete Stripe webhook endpoint (%s): %w", endpointID, err)
	}
	log.Printf("Stripe webhook endpoint deleted: ID=%s", endpointID)
	return nil
}
