package paymentx

import (
	"context"
	"fmt"
	"log"

	"github.com/braintree-go/braintree-go"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type braintreePaymentClient struct {
	gateway         *braintree.Braintree
	merchantAccount string
}

func newBraintreePaymentClient(connection *models.Connection, crypto cryptox.ICrypto) (IPaymentClient, error) {
	secret, err := crypto.Decrypt(connection.EncryptedKey)
	if err != nil {
		return nil, err
	}
	creds, err := ParseBraintreeCredentials(secret)
	if err != nil {
		return nil, err
	}
	env, err := braintree.EnvironmentFromName(creds.Environment)
	if err != nil {
		return nil, err
	}
	gateway := braintree.New(env, creds.MerchantID, connection.PublicKey, creds.PrivateKey)
	return &braintreePaymentClient{gateway: gateway, merchantAccount: creds.MerchantAccount}, nil
}

func (c *braintreePaymentClient) CreateCustomer(email, name string) (string, error) {
	customer, err := c.gateway.Customer().Create(context.Background(), &braintree.CustomerRequest{Email: email, FirstName: name})
	if err != nil {
		return "", fmt.Errorf("error creating Braintree customer: %w", err)
	}
	return customer.Id, nil
}

func (c *braintreePaymentClient) CreateCheckoutSession(customerID string) (string, string, error) {
	token, err := c.gateway.ClientToken().GenerateWithCustomer(context.Background(), customerID)
	if err != nil {
		return "", "", fmt.Errorf("error creating Braintree client token: %w", err)
	}
	return dbx.GenPublicID("btsess"), token, nil
}

func (c *braintreePaymentClient) CreateCharge(customerID string, total int64, currency string, paymentMethodID string, metadata map[string]string) error {
	txReq := &braintree.TransactionRequest{
		Type:               "sale",
		Amount:             braintree.NewDecimal(total, 2),
		CustomerID:         customerID,
		PaymentMethodToken: paymentMethodID,
		OrderId:            braintreeOrderID(metadata),
		Options:            &braintree.TransactionOptions{SubmitForSettlement: true},
	}
	if c.merchantAccount != "" {
		txReq.MerchantAccountId = c.merchantAccount
	}
	tx, err := c.gateway.Transaction().Create(context.Background(), txReq)
	if err != nil {
		return fmt.Errorf("failed to create braintree transaction for customer (%s): %w", customerID, err)
	}
	if tx == nil || (tx.Status != braintree.TransactionStatusSubmittedForSettlement && tx.Status != braintree.TransactionStatusSettling && tx.Status != braintree.TransactionStatusSettled) {
		return fmt.Errorf("braintree transaction for customer (%s) failed with status: %s", customerID, tx.Status)
	}
	log.Printf("Braintree transaction for customer (%s) created successfully with ID: %s", customerID, tx.Id)
	return nil
}

func (c *braintreePaymentClient) Verify() error {
	_, err := c.gateway.ClientToken().Generate(context.Background())
	if err != nil {
		return fmt.Errorf("failed to verify Braintree payment client: %w", err)
	}
	return nil
}

func (c *braintreePaymentClient) RegisterWebhook(webhookURL string) (*WebhookRegistration, error) {
	log.Printf("Braintree webhook requires manual setup. Configure webhook URL: %s", webhookURL)
	return &WebhookRegistration{}, nil
}

func (c *braintreePaymentClient) DeleteWebhook(endpointID string) error { return nil }
