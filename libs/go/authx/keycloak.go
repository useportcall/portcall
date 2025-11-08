package authx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type localAuthClient struct {
	jwksURL, issuer, audience string
	k                         keyfunc.Keyfunc
}

func newLocalAuthClient() IAuthClient {
	issuer := "http://localhost:8080/realms/dev"         // MUST match token iss exactly, TODO: fix
	jwksURL := issuer + "/protocol/openid-connect/certs" // TODO: fix
	// v3: build a keyfunc that auto-refreshes JWKS in the background.
	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		// handle this properly in prod (log, panic, etc.)
		panic(err)
	}

	return &localAuthClient{
		jwksURL: jwksURL,
		issuer:  issuer,
		// audience: "example",
		k: k,
	}
}

func (client *localAuthClient) Validate(ctx context.Context, header http.Header) (string, error) {
	h := header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		log.Print("Authorization header missing or invalid")
		return "", fmt.Errorf("authorization header missing or invalid")
	}
	raw := strings.TrimPrefix(h, "Bearer ")

	tok, err := jwt.Parse(raw, client.k.Keyfunc,
		jwt.WithIssuer(client.issuer),
		// jwt.WithAudience(audience),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithLeeway(2*time.Minute),
	)
	if err != nil || tok == nil || !tok.Valid {
		log.Printf("JWT parse error: %v", err)
		return "", fmt.Errorf("JWT parse error: %v", err)
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok {
		if email, _ := claims["email"].(string); email != "" {
			return email, nil
		}
	}

	return "", fmt.Errorf("email claim not found in token")
}
