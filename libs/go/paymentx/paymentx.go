package paymentx

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func New(connection *models.Connection, crypto cryptox.ICrypto) (IPaymentClient, error) {
	source := connection.Source

	switch source {
	case "local":
		return &localPaymentClient{}, nil
	case "stripe":
		stpSecretKey, err := crypto.Decrypt(connection.EncryptedKey)
		if err != nil {
			return nil, err
		}

		stp := &client.API{}
		stp.Init(stpSecretKey, nil)

		return stripePaymentClient{stp: stp}, nil
	default:
		return nil, fmt.Errorf("unsupported payment source: %s", source)
	}
}

type IPaymentClient interface {
	CreateCustomer(email, name string) (string, error)
	CreateCheckoutSession(customerID string) (string, string, error)
	CreateCharge(customerID string, total int64, currency string, paymentMethodID string) error
	Verify() error
}

// local

type localPaymentClient struct {
}

func (c *localPaymentClient) CreateCustomer(email, name string) (string, error) {
	id := dbx.GenPublicID("cust")
	return id, nil
}

func (c *localPaymentClient) CreateCheckoutSession(customerID string) (string, string, error) {
	id := dbx.GenPublicID("rand")
	clientSecret := dbx.GenPublicID("cls")
	return id, clientSecret, nil
}

func (c *localPaymentClient) CreateCharge(customerID string, total int64, currency string, paymentMethodID string) error {
	log.Printf("Local payment charge created for customer (%s) with amount: %d %s", customerID, total, currency)
	return nil
}

func (c *localPaymentClient) Verify() error {
	return nil
}

// stripe

type stripePaymentClient struct {
	stp *client.API
}

func (c stripePaymentClient) CreateCustomer(email, name string) (string, error) {
	customerParams := &stripe.CustomerParams{}
	customerParams.Email = stripe.String(email)
	customerParams.Name = stripe.String(name)

	stpCustomer, err := c.stp.Customers.New(customerParams)
	if err != nil {
		return "", fmt.Errorf("error creating Stripe customer: %w", err)
	}

	return stpCustomer.ID, nil
}

func (c stripePaymentClient) CreateCheckoutSession(customerID string) (string, string, error) {
	stpSetupIntentParams := &stripe.SetupIntentParams{}
	stpSetupIntentParams.Customer = stripe.String(customerID)
	stpSetupIntentParams.PaymentMethodTypes = stripe.StringSlice([]string{"card"})
	stpSetupIntent, err := c.stp.SetupIntents.New(stpSetupIntentParams)
	if err != nil {
		return "", "", err
	}
	return stpSetupIntent.ID, stpSetupIntent.ClientSecret, nil
}

func (c stripePaymentClient) CreateCharge(customerID string, total int64, currency string, paymentMethodID string) error {
	params := &stripe.PaymentIntentParams{
		Customer:           stripe.String(customerID),
		Amount:             stripe.Int64(total),
		Currency:           stripe.String(currency),
		PaymentMethod:      stripe.String(paymentMethodID),
		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodManual)),
		Confirm:            stripe.Bool(true),
	}
	stpPaymentIntent, err := c.stp.PaymentIntents.New(params)
	if err != nil {
		return fmt.Errorf("failed to create payment intent for customer (%s): %w", customerID, err)
	}

	if stpPaymentIntent.Status == stripe.PaymentIntentStatusRequiresAction {
		return fmt.Errorf("payment intent for customer (%s) requires action: %w", customerID, err)
	}

	if stpPaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
		return fmt.Errorf("payment intent for customer (%s) status is not succeeded: %s", customerID, stpPaymentIntent.Status)
	}

	log.Printf("Payment intent for customer (%s) created successfully with ID: %s", customerID, stpPaymentIntent.ID)
	return nil
}

func (c stripePaymentClient) Verify() error {
	_, err := c.stp.Account.Get()
	if err != nil {
		return fmt.Errorf("failed to verify Stripe payment client: %w", err)
	}
	return nil
}
