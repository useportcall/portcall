package authx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type localAuthClient struct {
	jwksURL, issuer, audience string
	k                         keyfunc.Keyfunc
}

func newLocalAuthClient() (IAuthClient, error) {
	apiURL := os.Getenv("KEYCLOAK_API_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("KEYCLOAK_API_URL environment variable is not set")
	}

	// Internal URL for fetching JWKS (accessible from within Kubernetes)
	internalURL := os.Getenv("KEYCLOAK_INTERNAL_URL")
	if internalURL == "" {
		internalURL = apiURL // fallback to external URL if internal not set
	}

	// Get realm from environment variable, default to "dev" for backwards compatibility
	realm := os.Getenv("KEYCLOAK_REALM")
	if realm == "" {
		realm = "dev"
	}

	// Get client ID for audience validation
	// Defaults based on realm: "portcall-admin" for admin realm, "portcall" otherwise
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	if clientID == "" {
		if realm == "admin" {
			clientID = "portcall-admin"
		} else {
			clientID = "portcall"
		}
	}

	// External URL for issuer validation (matches token issuer)
	issuer := apiURL + "/realms/" + realm
	jwksURL := internalURL + "/realms/" + realm + "/protocol/openid-connect/certs"

	log.Printf("Keycloak auth client initialized - JWKS URL: %s, Expected Issuer: %s, Expected Audience: %s", jwksURL, issuer, clientID)

	// v3: build a keyfunc that auto-refreshes JWKS in the background.
	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS keyfunc: %w", err)
	}

	return &localAuthClient{
		jwksURL:  jwksURL,
		issuer:   issuer,
		audience: clientID,
		k:        k,
	}, nil
}

func (client *localAuthClient) Validate(ctx context.Context, header http.Header) (*Claims, error) {
	h := header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		log.Print("Authorization header missing or invalid")
		return nil, fmt.Errorf("authorization header missing or invalid")
	}
	raw := strings.TrimPrefix(h, "Bearer ")

	tok, err := jwt.Parse(raw, client.k.Keyfunc,
		jwt.WithIssuer(client.issuer),
		jwt.WithAudience(client.audience),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithLeeway(2*time.Minute),
	)
	if err != nil || tok == nil || !tok.Valid {
		if tok != nil {
			if claims, ok := tok.Claims.(jwt.MapClaims); ok {
				log.Printf("JWT parse error: %v (Expected issuer: %s, Token issuer: %v)", err, client.issuer, claims["iss"])
			} else {
				log.Printf("JWT parse error: %v", err)
			}
		} else {
			log.Printf("JWT parse error: %v", err)
		}
		return nil, fmt.Errorf("JWT parse error: %v", err)
	}

	if claims, ok := tok.Claims.(jwt.MapClaims); ok {
		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("email claim not found in token")
		}

		givenName := parseClaimString(claims, "given_name")
		familyName := parseClaimString(claims, "family_name")

		// Fallback to "name" when split claims are missing from token/client scopes.
		if givenName == nil {
			givenName, familyName = parseFullNameClaim(claims)
		}

		response := new(Claims)
		response.Email = email
		response.GivenName = givenName
		response.FamilyName = familyName

		return response, nil
	}

	return nil, fmt.Errorf("email claim not found in token")
}

func parseClaimString(claims jwt.MapClaims, key string) *string {
	v, ok := claims[key]
	if !ok {
		return nil
	}

	s, ok := v.(string)
	if !ok {
		return nil
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	return &s
}

func parseFullNameClaim(claims jwt.MapClaims) (*string, *string) {
	name := parseClaimString(claims, "name")
	if name == nil {
		return nil, nil
	}

	parts := strings.Fields(*name)
	if len(parts) == 0 {
		return nil, nil
	}

	givenName := parts[0]
	if len(parts) == 1 {
		return &givenName, nil
	}

	familyName := strings.Join(parts[1:], " ")
	return &givenName, &familyName
}
