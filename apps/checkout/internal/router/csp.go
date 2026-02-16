package router

import (
	"os"
	"strings"
)

func buildCheckoutCSP() string {
	var directives []string
	braintreeConnect := " https://api.braintreegateway.com https://client-analytics.braintreegateway.com https://*.braintree-api.com https://api.sandbox.braintreegateway.com https://client-analytics.sandbox.braintreegateway.com"
	connectSrc := "connect-src 'self' https://api.stripe.com https://r.stripe.com" + braintreeConnect
	if os.Getenv("APP_ENV") != "production" {
		connectSrc = "connect-src 'self' https://api.stripe.com https://r.stripe.com" + braintreeConnect + " ws: wss:"
	}

	directives = append(directives, "default-src 'self'")
	directives = append(directives, "base-uri 'self'")
	directives = append(directives, "object-src 'none'")
	directives = append(directives, "frame-ancestors 'none'")
	directives = append(directives, "img-src 'self' data: https: https://assets.braintreegateway.com https://*.paypalobjects.com")
	directives = append(directives, "style-src 'self' 'unsafe-inline'")
	directives = append(directives, "font-src 'self' data:")
	directives = append(directives, "script-src 'self' 'unsafe-inline' https://js.stripe.com https://js.braintreegateway.com https://assets.braintreegateway.com")
	directives = append(directives, connectSrc)
	directives = append(directives, "frame-src 'self' https://js.stripe.com https://hooks.stripe.com https://assets.braintreegateway.com https://*.paypalobjects.com https://c.paypal.com")
	directives = append(directives, "form-action 'self'")

	return strings.Join(directives, "; ")
}
