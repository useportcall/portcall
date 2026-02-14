package checkout_session

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
)

// Create handles checkout session creation.
func (s *service) Create(input *CreateInput) (*CreateResult, error) {
	redirectURL, err := validateCheckoutReturnURL(input.RedirectURL, "redirect_url")
	if err != nil {
		return nil, err
	}
	cancelURL, err := validateCheckoutReturnURL(input.CancelURL, "cancel_url")
	if err != nil {
		return nil, err
	}

	var company models.Company
	if err := s.db.FindFirst(&company, "app_id = ?", input.AppID); err != nil {
		if input.RequireBillingAddress {
			return nil, fmt.Errorf("error retrieving company for app: %w", err)
		}
		log.Println("no company found for app, continuing without billing address")
	}
	if input.RequireBillingAddress && company.BillingAddressID == 0 {
		return nil, NewValidationError("company does not have a billing address")
	}

	var plan models.Plan
	if err := s.db.GetForPublicID(input.AppID, input.PlanID, &plan); err != nil {
		return nil, err
	}
	if plan.Status != "published" {
		return nil, NewValidationError("plan with id '%s' is not yet published", plan.PublicID)
	}

	var config models.AppConfig
	if err := s.db.FindFirst(&config, "app_id = ?", input.AppID); err != nil {
		return nil, err
	}

	connection, err := s.resolveConnection(input.AppID, config.DefaultConnectionID)
	if err != nil {
		return nil, err
	}

	payment, err := paymentx.New(connection, s.crypto)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.GetForPublicID(input.AppID, input.UserID, &user); err != nil {
		return nil, err
	}

	if user.PaymentCustomerID == "" {
		customerID, err := payment.CreateCustomer(user.Email, user.Name)
		if err != nil {
			return nil, err
		}
		user.PaymentCustomerID = customerID
		if err := s.db.Save(&user); err != nil {
			return nil, err
		}
	}

	sessionID, clientSecret, err := payment.CreateCheckoutSession(user.PaymentCustomerID)
	if err != nil {
		return nil, err
	}

	var companyAddressID *uint
	if company.BillingAddressID != 0 {
		companyAddressID = &company.BillingAddressID
	}

	session := &models.CheckoutSession{
		PublicID:             dbx.GenPublicID("cs"),
		ExpiresAt:            time.Now().Add(48 * time.Hour),
		AppID:                input.AppID,
		UserID:               user.ID,
		PlanID:               plan.ID,
		ExternalClientSecret: clientSecret,
		ExternalPublicKey:    connection.PublicKey,
		ExternalSessionID:    sessionID,
		ExternalProvider:     connection.Source,
		CancelURL:            &cancelURL,
		RedirectURL:          &redirectURL,
		CompanyAddressID:     companyAddressID,
	}

	if err := s.db.Create(session); err != nil {
		return nil, err
	}

	session.User = user
	session.Plan = plan

	checkoutURL, err := buildCheckoutURL(os.Getenv("CHECKOUT_URL"), session, s.crypto)
	if err != nil {
		return nil, err
	}

	return &CreateResult{
		Session:     session,
		CheckoutURL: checkoutURL,
	}, nil
}

func validateCheckoutReturnURL(raw string, field string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", NewValidationError("invalid %s", field)
	}

	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return "", NewValidationError("invalid %s", field)
	}
	if parsed.Host == "" || parsed.User != nil {
		return "", NewValidationError("invalid %s", field)
	}

	return parsed.String(), nil
}

func buildCheckoutURL(baseURL string, session *models.CheckoutSession, crypto cryptox.ICrypto) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" || session == nil {
		return "", nil
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid CHECKOUT_URL: %w", err)
	}

	query := parsed.Query()
	query.Set("id", session.PublicID)

	if crypto != nil {
		token, err := cryptox.CreateCheckoutSessionToken(crypto, session.PublicID, session.ExpiresAt)
		if err != nil {
			return "", err
		}
		query.Set("st", token)
	}

	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
