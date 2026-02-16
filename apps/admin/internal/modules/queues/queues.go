package queues

import (
	"github.com/useportcall/portcall/libs/go/routerx"
)

type QueueInfo struct {
	Name  string     `json:"name"`
	Tasks []TaskInfo `json:"tasks"`
}

type TaskInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Params      []ParamInfo `json:"params"`
}

type ParamInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type EnqueueRequest struct {
	QueueName string                 `json:"queue_name"`
	TaskType  string                 `json:"task_type"`
	Payload   map[string]interface{} `json:"payload"`
}

// GetQueues returns all available queues and their tasks
func GetQueues(c *routerx.Context) {
	queues := []QueueInfo{
		{
			Name: "billing_queue",
			Tasks: []TaskInfo{
				{
					Name:        "create_payment_method",
					Description: "Create a payment method from a Stripe session",
					Params: []ParamInfo{
						{Name: "external_session_id", Type: "string", Required: true},
						{Name: "external_payment_method_id", Type: "string", Required: true},
					},
				},
				{
					Name:        "create_subscription",
					Description: "Create a new subscription for a user",
					Params: []ParamInfo{
						{Name: "plan_id", Type: "int", Required: true},
						{Name: "user_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "update_subscription",
					Description: "Update an existing subscription to a new plan",
					Params: []ParamInfo{
						{Name: "subscription_id", Type: "int", Required: true},
						{Name: "plan_id", Type: "int", Required: true},
						{Name: "app_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "create_subscription_items",
					Description: "Create subscription items for a subscription",
					Params: []ParamInfo{
						{Name: "subscription_id", Type: "int", Required: true},
						{Name: "plan_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "create_invoice",
					Description: "Create an invoice for a subscription",
					Params: []ParamInfo{
						{Name: "subscription_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "create_upgrade_invoice",
					Description: "Create an invoice for a subscription upgrade with prorated amount",
					Params: []ParamInfo{
						{Name: "subscription_id", Type: "int", Required: true},
						{Name: "price_difference", Type: "int", Required: true},
						{Name: "old_plan_id", Type: "int", Required: true},
						{Name: "new_plan_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "create_invoice_items",
					Description: "Create line items for an invoice",
					Params: []ParamInfo{
						{Name: "invoice_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "calculate_invoice_totals",
					Description: "Calculate and update invoice totals",
					Params: []ParamInfo{
						{Name: "invoice_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "pay_invoice",
					Description: "Process payment for an invoice",
					Params: []ParamInfo{
						{Name: "invoice_id", Type: "int", Required: true},
						{Name: "skip_next", Type: "bool", Required: false},
					},
				},
				{
					Name:        "resolve_invoice",
					Description: "Resolve and finalize an invoice after payment",
					Params: []ParamInfo{
						{Name: "invoice_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "process_stripe_payment_failure",
					Description: "Process Stripe webhook payment failures for dunning",
					Params: []ParamInfo{
						{Name: "invoice_id", Type: "int", Required: true},
						{Name: "attempt", Type: "int", Required: false},
						{Name: "no_retry", Type: "bool", Required: false},
						{Name: "event_type", Type: "string", Required: false},
						{Name: "failure_reason", Type: "string", Required: false},
					},
				},
				{
					Name:        "process_braintree_webhook_event",
					Description: "Process Braintree webhook events for checkout resolution and dunning",
					Params: []ParamInfo{
						{Name: "kind", Type: "string", Required: true},
						{Name: "order_id", Type: "string", Required: false},
						{Name: "payment_method_token", Type: "string", Required: false},
						{Name: "failure_count", Type: "int", Required: false},
						{Name: "failure_reason", Type: "string", Required: false},
					},
				},
				{
					Name:        "find_subscriptions_to_reset",
					Description: "Find subscriptions that need to be reset based on billing interval",
					Params:      []ParamInfo{},
				},
				{
					Name:        "start_subscription_reset",
					Description: "Start the reset process for a subscription",
					Params: []ParamInfo{
						{Name: "subscription_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "end_subscription_reset",
					Description: "Complete the reset process for a subscription",
					Params: []ParamInfo{
						{Name: "subscription_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "process_meter_event",
					Description: "Process a metered usage event",
					Params: []ParamInfo{
						{Name: "meter_event_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "create_entitlements",
					Description: "Create entitlements for a user based on their plan",
					Params: []ParamInfo{
						{Name: "user_id", Type: "int", Required: true},
						{Name: "plan_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "process_plan_switch",
					Description: "Process a plan switch for a subscription",
					Params: []ParamInfo{
						{Name: "old_plan_id", Type: "int", Required: true},
						{Name: "new_plan_id", Type: "int", Required: true},
						{Name: "subscription_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "df_increment",
					Description: "Dogfood: Increment a feature usage counter",
					Params: []ParamInfo{
						{Name: "user_id", Type: "string", Required: true},
						{Name: "feature", Type: "string", Required: true},
						{Name: "is_test", Type: "bool", Required: true},
					},
				},
				{
					Name:        "df_decrement",
					Description: "Dogfood: Decrement a feature usage counter",
					Params: []ParamInfo{
						{Name: "user_id", Type: "string", Required: true},
						{Name: "feature", Type: "string", Required: true},
						{Name: "is_test", Type: "bool", Required: true},
					},
				},
				{
					Name:        "df_create_user",
					Description: "Dogfood: Create a test user for an app",
					Params: []ParamInfo{
						{Name: "app_id", Type: "int", Required: true},
					},
				},
				{
					Name:        "df_create_subscription",
					Description: "Dogfood: Create a test subscription for an app",
					Params: []ParamInfo{
						{Name: "app_id", Type: "int", Required: true},
					},
				},
			},
		},
		{
			Name: "email_queue",
			Tasks: []TaskInfo{
				{
					Name:        "send_invoice_paid_email",
					Description: "Send invoice paid confirmation email",
					Params: []ParamInfo{
						{Name: "recipient_email", Type: "string", Required: true},
						{Name: "recipient_name", Type: "string", Required: false},
						{Name: "invoice_number", Type: "string", Required: false},
						{Name: "company_name", Type: "string", Required: false},
						{Name: "amount_paid", Type: "string", Required: false},
						{Name: "date_paid", Type: "string", Required: false},
					},
				},
				{
					Name:        "send_invoice_dunning_email",
					Description: "Send payment failed/dunning notification email",
					Params: []ParamInfo{
						{Name: "recipient_email", Type: "string", Required: true},
						{Name: "invoice_number", Type: "string", Required: true},
						{Name: "amount_due", Type: "string", Required: false},
						{Name: "due_date", Type: "string", Required: false},
						{Name: "attempt", Type: "int", Required: false},
						{Name: "max_attempts", Type: "int", Required: false},
					},
				},
				{
					Name:        "send_quote_email",
					Description: "Send quote to recipient",
					Params: []ParamInfo{
						{Name: "recipient_email", Type: "string", Required: true},
						{Name: "CustomerName", Type: "string", Required: false},
						{Name: "CompanyName", Type: "string", Required: false},
						{Name: "QuoteURL", Type: "string", Required: false},
					},
				},
				{
					Name:        "send_quote_accepted_confirmation_email",
					Description: "Send confirmation email when a quote is accepted",
					Params: []ParamInfo{
						{Name: "recipient_email", Type: "string", Required: true},
						{Name: "RecipientName", Type: "string", Required: false},
						{Name: "QuoteLink", Type: "string", Required: false},
					},
				},
				{
					Name:        "process_postmark_webhook_event",
					Description: "Process a Postmark webhook event (bounce, spam complaint, etc.)",
					Params: []ParamInfo{
						{Name: "raw_event", Type: "object", Required: true},
					},
				},
			},
		},
	}

	c.OK(queues)
}

// EnqueueTask enqueues a task to the specified queue
func EnqueueTask(c *routerx.Context) {
	var req EnqueueRequest
	if err := c.BindJSON(&req); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if req.QueueName == "" || req.TaskType == "" {
		c.BadRequest("queue_name and task_type are required")
		return
	}

	err := c.Queue().Enqueue(req.TaskType, req.Payload, req.QueueName)
	if err != nil {
		c.ServerError("Failed to enqueue task", err)
		return
	}

	c.OK(map[string]interface{}{
		"success":    true,
		"message":    "Task '" + req.TaskType + "' enqueued to '" + req.QueueName + "'",
		"queue_name": req.QueueName,
		"task_type":  req.TaskType,
	})
}
