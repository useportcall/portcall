package authx

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

type IAuthClient interface {
	Validate(ctx context.Context, header http.Header) (string, error)
}

type clerkAuthClient struct {
	config     *clerk.ClientConfig
	store      JWKStore
	jwksClient *jwks.Client
	user       *user.Client
}

func New() IAuthClient {
	authModule := os.Getenv("AUTH_MODULE")
	switch authModule { // TODO: explore build tags to exclude unused auth modules
	case "clerk":
		return newClerkAuthClient()
	case "local":
		return newLocalAuthClient()
	default:
		log.Fatal("AUTH_MODULE_NOT_FOUND")
		return nil
	}
}

func newClerkAuthClient() IAuthClient {
	clerkSecret := os.Getenv("CLERK_SECRET")
	if clerkSecret == "" {
		log.Fatal("CLERK_SECRET_NULL")
	}

	config := &clerk.ClientConfig{}
	config.Key = clerk.String(clerkSecret)

	store := NewJWKStore()
	jwksClient := jwks.NewClient(config)
	user := user.NewClient(config)

	return &clerkAuthClient{config: config, store: store, jwksClient: jwksClient, user: user}
}

func (client *clerkAuthClient) Validate(ctx context.Context, header http.Header) (string, error) {
	sessionToken := strings.TrimPrefix(header.Get("Authorization"), "Bearer ")

	// Attempt to get the JSON Web Key from your store.
	jwk := client.store.GetJWK()

	if jwk == nil {
		// Decode the session JWT so that we can find the key ID.
		unsafeClaims, err := jwt.Decode(ctx, &jwt.DecodeParams{
			Token: sessionToken,
		})
		if err != nil {
			log.Println("DECODE_ERR", err)
			return "", err
		}

		// Fetch the JSON Web Key
		jwk, err = jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
			KeyID:      unsafeClaims.KeyID,
			JWKSClient: client.jwksClient,
		})
		if err != nil {
			log.Println("JWK_ERR", err)
			return "", err
		}

		// Write the JSON Web Key to your store, so that next time
		// you can use the cached value.
		client.store.SetJWK(jwk)
	}

	// Verify the session
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: sessionToken,
		JWK:   jwk,
	})
	if err != nil {
		log.Println("CLAIMS_ERR", err)
		return "", err
	}

	usr, err := client.user.Get(ctx, claims.Subject)
	if err != nil {
		return "", err
	}

	return usr.EmailAddresses[0].EmailAddress, nil
}
