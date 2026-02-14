package webhookx

import (
	"fmt"

	"github.com/braintree-go/braintree-go"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func braintreeGatewayForConnection(c *routerx.Context, connectionID string) (*braintree.Braintree, error) {
	var connection models.Connection
	if err := c.DB().FindFirst(&connection, "public_id = ?", connectionID); err != nil {
		return nil, fmt.Errorf("connection not found")
	}
	if connection.Source != "braintree" {
		return nil, fmt.Errorf("connection source is not braintree")
	}
	secret, err := c.Crypto().Decrypt(connection.EncryptedKey)
	if err != nil {
		return nil, err
	}
	creds, err := paymentx.ParseBraintreeCredentials(secret)
	if err != nil {
		return nil, err
	}
	env, err := braintree.EnvironmentFromName(creds.Environment)
	if err != nil {
		return nil, err
	}
	return braintree.New(env, creds.MerchantID, connection.PublicKey, creds.PrivateKey), nil
}
