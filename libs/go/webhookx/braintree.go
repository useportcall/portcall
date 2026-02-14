package webhookx

import (
	"log"
	"net/http"
	"strconv"

	"github.com/useportcall/portcall/libs/go/routerx"
)

func (w *Router) HandleBraintreeChallenge(c *routerx.Context) {
	if !w.guard(c, "braintree") {
		return
	}
	challenge := c.Query("bt_challenge")
	connectionID, ok := c.Params.Get("connection_id")
	if !ok || challenge == "" {
		respondGenericError(c, http.StatusBadRequest)
		return
	}
	gateway, err := braintreeGatewayForConnection(c, connectionID)
	if err != nil {
		log.Printf("[webhookx] braintree challenge config error: %v", err)
		respondGenericError(c, http.StatusBadRequest)
		return
	}
	response, err := gateway.WebhookNotification().Verify(challenge)
	if err != nil {
		log.Printf("[webhookx] braintree challenge verify failed: %v", err)
		respondGenericError(c, http.StatusBadRequest)
		return
	}
	c.String(http.StatusOK, response)
}

func (w *Router) HandleBraintreeWebhook(c *routerx.Context) {
	if !w.guard(c, "braintree") {
		return
	}
	connectionID, ok := c.Params.Get("connection_id")
	if !ok {
		respondBraintreeOK(c)
		return
	}
	gateway, err := braintreeGatewayForConnection(c, connectionID)
	if err != nil {
		log.Printf("[webhookx] braintree config error: %v", err)
		respondBraintreeOK(c)
		return
	}
	if err := c.Request.ParseForm(); err != nil {
		log.Printf("[webhookx] braintree form parse failed: %v", err)
		respondBraintreeOK(c)
		return
	}
	signature := c.Request.PostFormValue("bt_signature")
	payload := c.Request.PostFormValue("bt_payload")
	if signature == "" || payload == "" {
		respondBraintreeOK(c)
		return
	}
	notification, err := gateway.WebhookNotification().Parse(signature, payload)
	if err != nil || !relevantBraintreeEvents[notification.Kind] {
		if err != nil {
			log.Printf("[webhookx] braintree signature parse failed: %v", err)
		}
		respondBraintreeOK(c)
		return
	}
	job := map[string]any{"kind": notification.Kind}
	if notification.Subject != nil && notification.Subject.Transaction != nil {
		tx := notification.Subject.Transaction
		job["order_id"] = tx.OrderId
		job["payment_method_token"] = tx.PaymentMethodToken
		job["failure_reason"] = tx.ProcessorResponseText
	}
	if notification.Subject != nil && notification.Subject.Subscription != nil {
		sub := notification.Subject.Subscription
		if sub.PaymentMethodToken != "" {
			job["payment_method_token"] = sub.PaymentMethodToken
		}
		if attempt, err := strconv.Atoi(sub.FailureCount); err == nil && attempt > 0 {
			job["failure_count"] = attempt
		}
	}
	if err := c.Queue().Enqueue("process_braintree_webhook_event", job, "billing_queue"); err != nil {
		log.Printf("[webhookx] braintree enqueue failed: %v", err)
		respondGenericError(c, http.StatusInternalServerError)
		return
	}
	respondBraintreeOK(c)
}
