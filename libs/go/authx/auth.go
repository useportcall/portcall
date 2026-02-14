package authx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

type Claims struct {
	Email      string  `json:"email"`
	GivenName  *string `json:"given_name"`
	FamilyName *string `json:"family_name"`
}

type IAuthClient interface {
	Validate(ctx context.Context, header http.Header) (*Claims, error)
}

type clerkAuthClient struct {
	config     *clerk.ClientConfig
	store      JWKStore
	jwksClient *jwks.Client
	user       *user.Client
}

func New() (IAuthClient, error) { return NewFromEnv() }

func NewFromEnv() (IAuthClient, error) {
	authModule := os.Getenv("AUTH_MODULE")
	switch authModule { // TODO: explore build tags to exclude unused auth modules
	case "clerk":
		return newClerkAuthClient()
	default:
		return newLocalAuthClient()
	}
}

func newClerkAuthClient() (IAuthClient, error) {
	clerkSecret := os.Getenv("CLERK_SECRET")
	if clerkSecret == "" {
		return nil, fmt.Errorf("CLERK_SECRET environment variable is not set")
	}

	config := &clerk.ClientConfig{}
	config.Key = clerk.String(clerkSecret)

	store := NewJWKStore()
	jwksClient := jwks.NewClient(config)
	user := user.NewClient(config)

	return &clerkAuthClient{config: config, store: store, jwksClient: jwksClient, user: user}, nil
}

func (client *clerkAuthClient) Validate(ctx context.Context, header http.Header) (*Claims, error) {
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
			return nil, err
		}

		// Fetch the JSON Web Key
		jwk, err = jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
			KeyID:      unsafeClaims.KeyID,
			JWKSClient: client.jwksClient,
		})
		if err != nil {
			log.Println("JWK_ERR", err)
			return nil, err
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
		return nil, err
	}

	usr, err := client.user.Get(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	response := new(Claims)
	response.Email = usr.EmailAddresses[0].EmailAddress
	response.GivenName = usr.FirstName
	response.FamilyName = usr.LastName

	return response, nil
}
