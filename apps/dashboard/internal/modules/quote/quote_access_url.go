package quote

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func buildQuoteAccessURL(baseURL string, quote *models.Quote, crypto cryptox.ICrypto, now time.Time) (*string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return nil, err
	}

	path := strings.TrimRight(parsed.Path, "/")
	parsed.Path = path + fmt.Sprintf("/quotes/%s", quote.PublicID)

	expiry := quoteExpiry(quote, now)
	token, err := cryptox.CreateQuoteAccessToken(crypto, quote.PublicID, expiry)
	if err != nil {
		return nil, err
	}

	query := parsed.Query()
	query.Set("qt", token)
	parsed.RawQuery = query.Encode()

	value := parsed.String()
	return &value, nil
}
