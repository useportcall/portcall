package router

import (
	"os"
	"strings"
)

func buildCheckoutCSP() string {
	var directives []string
	connectSrc := "connect-src 'self' https://api.stripe.com https://r.stripe.com"
	if os.Getenv("APP_ENV") != "production" {
		connectSrc = "connect-src 'self' https://api.stripe.com https://r.stripe.com ws: wss:"
	}

	directives = append(directives, "default-src 'self'")
	directives = append(directives, "base-uri 'self'")
	directives = append(directives, "object-src 'none'")
	directives = append(directives, "frame-ancestors 'none'")
	directives = append(directives, "img-src 'self' data: https:")
	directives = append(directives, "style-src 'self' 'unsafe-inline'")
	directives = append(directives, "font-src 'self' data:")
	directives = append(directives, "script-src 'self' 'unsafe-inline' https://js.stripe.com")
	directives = append(directives, connectSrc)
	directives = append(directives, "frame-src 'self' https://js.stripe.com https://hooks.stripe.com")
	directives = append(directives, "form-action 'self'")

	return strings.Join(directives, "; ")
}
