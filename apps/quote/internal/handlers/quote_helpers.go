package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func convertCentsToDollars(cents int64) string {
	return fmt.Sprintf("$%.2f", float64(cents)/100.0)
}

func toBasePrice(cents int64, interval string) string {
	return fmt.Sprintf("$%.2f / %s", float64(cents)/100.0, interval)
}

func convertSnakeToCamelCaps(snakeStr string) string {
	parts := strings.Split(snakeStr, "_")
	for i, part := range parts {
		parts[i] = capitalizeFirstLetter(part)
	}
	return strings.Join(parts, " ")
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

func markQuoteAccepted(quote *models.Quote, now time.Time) {
	quote.Status = "accepted"
	quote.AcceptedAt = &now
	sigURL := fmt.Sprintf("s3://quote-signatures/%s.png", quote.PublicID)
	quote.SignatureURL = &sigURL
}

func validateSignature(c *routerx.Context) ([]byte, bool) {
	signatureData := c.PostForm("signatureData")
	if signatureData == "" {
		c.String(http.StatusBadRequest, "No signature data provided")
		return nil, false
	}
	if len(signatureData) > 2*1024*1024 {
		c.String(http.StatusBadRequest, "Signature data too large")
		return nil, false
	}
	imgBytes, err := decodeSignatureData(signatureData)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid signature data")
		return nil, false
	}
	if !hasSignatureStroke(imgBytes) {
		c.String(http.StatusBadRequest, "Signature cannot be empty")
		return nil, false
	}
	return imgBytes, true
}

func buildQuoteLink(crypto cryptox.ICrypto, quoteAppURL string, quote *models.Quote, now time.Time) (string, error) {
	tokenExpiry := now.Add(90 * 24 * time.Hour)
	if quote.ExpiresAt != nil && quote.ExpiresAt.After(tokenExpiry) {
		tokenExpiry = quote.ExpiresAt.UTC()
	}
	token, err := cryptox.CreateQuoteAccessToken(crypto, quote.PublicID, tokenExpiry)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/quotes/%s?qt=%s", strings.TrimRight(quoteAppURL, "/"), quote.PublicID, url.QueryEscape(token)), nil
}
