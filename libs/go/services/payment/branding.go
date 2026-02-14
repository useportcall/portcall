package payment

import (
	"os"
	"strings"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func emailLogoURL(company *models.Company) string {
	if company != nil && strings.TrimSpace(company.IconLogoURL) != "" {
		return company.IconLogoURL
	}
	return strings.TrimSpace(os.Getenv("EMAIL_BRAND_LOGO_URL"))
}
