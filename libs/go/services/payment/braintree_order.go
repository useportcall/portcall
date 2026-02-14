package payment

import (
	"strconv"
	"strings"
)

func parseOrderMetadataUint(orderID string, key string) uint {
	value := parseOrderMetadataValue(orderID, key)
	if value == "" {
		return 0
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}
	return uint(parsed)
}

func parseOrderMetadataValue(orderID string, key string) string {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return ""
	}
	if value := parseKeyedOrderValue(orderID, key); value != "" {
		return value
	}
	if key == "portcall_invoice_id" && strings.HasPrefix(orderID, "portcall_invoice_") {
		return strings.TrimPrefix(orderID, "portcall_invoice_")
	}
	if key == "portcall_checkout_session_id" && strings.HasPrefix(orderID, "portcall_checkout_") {
		return strings.TrimPrefix(orderID, "portcall_checkout_")
	}
	return ""
}

func parseKeyedOrderValue(orderID string, key string) string {
	for _, token := range splitOrderTokens(orderID) {
		for _, separator := range []string{"=", ":"} {
			prefix := key + separator
			if strings.HasPrefix(token, prefix) {
				return strings.TrimSpace(strings.TrimPrefix(token, prefix))
			}
		}
	}
	return ""
}

func splitOrderTokens(orderID string) []string {
	return strings.FieldsFunc(orderID, func(r rune) bool {
		return r == '|' || r == ';' || r == ','
	})
}
