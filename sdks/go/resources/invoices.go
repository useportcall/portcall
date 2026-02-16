package resources

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Invoice represents an invoice in the system
type Invoice struct {
	ID             string     `json:"id"`
	SubscriptionID *string    `json:"subscription_id,omitempty"`
	InvoiceNumber  string     `json:"invoice_number"`
	Currency       string     `json:"currency"`
	Total          int64      `json:"total"`
	SubTotal       int64      `json:"subtotal"`
	Status         string     `json:"status"`
	DueBy          *time.Time `json:"due_by,omitempty"`
	PDFURL         *string    `json:"pdf_url,omitempty"`
	CustomerEmail  string     `json:"recipient_email"`
	CustomerID     string     `json:"recipient_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ListInvoicesParams are the parameters for listing invoices
type ListInvoicesParams struct {
	UserID         string
	SubscriptionID string
	Limit          int
}

// Invoices provides access to invoice-related API operations
type Invoices struct {
	http *HTTPClient
}

// NewInvoices creates a new Invoices resource
func NewInvoices(http *HTTPClient) *Invoices {
	return &Invoices{http: http}
}

// List returns all invoices matching the given parameters
func (i *Invoices) List(ctx context.Context, params *ListInvoicesParams) ([]Invoice, error) {
	query := url.Values{}
	if params != nil {
		if params.UserID != "" {
			query.Set("user_id", params.UserID)
		}
		if params.SubscriptionID != "" {
			query.Set("subscription_id", params.SubscriptionID)
		}
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
	}

	path := "/v1/invoices"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var resp DataWrapper[[]Invoice]
	if err := i.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns an invoice by ID
func (i *Invoices) Get(ctx context.Context, invoiceID string) (*Invoice, error) {
	var resp DataWrapper[Invoice]
	if err := i.http.Get(ctx, fmt.Sprintf("/v1/invoices/%s", invoiceID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
