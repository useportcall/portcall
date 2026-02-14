package paymentx

import "strconv"

func braintreeOrderID(metadata map[string]string) string {
	if metadata == nil {
		return ""
	}
	if invoiceID := metadata["portcall_invoice_id"]; invoiceID != "" {
		return "portcall_invoice_id=" + invoiceID
	}
	if sessionID := metadata["portcall_checkout_session_id"]; sessionID != "" {
		return "portcall_checkout_session_id=" + sessionID
	}
	if raw := metadata["invoice_id"]; raw != "" {
		if _, err := strconv.Atoi(raw); err == nil {
			return "portcall_invoice_id=" + raw
		}
	}
	return ""
}
