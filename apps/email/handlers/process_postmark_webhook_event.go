package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// ProcessPostmarkWebhookEvent processes incoming Postmark webhook events.
func ProcessPostmarkWebhookEvent(c server.IContext) error {
	var jobArgs map[string]any
	if err := json.Unmarshal(c.Payload(), &jobArgs); err != nil {
		return fmt.Errorf("failed to unmarshal job payload: %w", err)
	}

	rawEventInterface, ok := jobArgs["raw_event"]
	if !ok {
		return fmt.Errorf("missing raw_event in job args")
	}

	// rawEventInterface will be a map[string]interface{} or similar because of json.Unmarshal
	// We need to marshal it back to bytes to unmarshal into specific structs,
	// or use json.RawMessage in the map if we defined a struct for jobArgs.
	// Since we used map[string]any, it's likely a map now.

	rawBytes, err := json.Marshal(rawEventInterface)
	if err != nil {
		return fmt.Errorf("failed to marshal raw event: %w", err)
	}

	// First unmarshal into base to get the type
	var base emailx.PostmarkBaseEvent
	if err := json.Unmarshal(rawBytes, &base); err != nil {
		return fmt.Errorf("failed to unmarshal base event: %w", err)
	}

	log.Printf("Processing Postmark event: %s (MessageID: %s)", base.RecordType, base.MessageID)

	switch base.RecordType {
	case emailx.PostmarkTypeBounce:
		var bounce emailx.PostmarkBounceEvent
		if err := json.Unmarshal(rawBytes, &bounce); err != nil {
			return fmt.Errorf("failed to unmarshal bounce event: %w", err)
		}
		handleBounce(bounce)

	case emailx.PostmarkTypeSpamComplaint:
		var spam emailx.PostmarkSpamComplaintEvent
		if err := json.Unmarshal(rawBytes, &spam); err != nil {
			return fmt.Errorf("failed to unmarshal spam event: %w", err)
		}
		handleSpamComplaint(spam)

	case "SubscriptionChange":
		// Handle subscription change if implemented
		log.Printf("Received SubscriptionChange event from Postmark")

	default:
		log.Printf("Unhandled Postmark event type: %s", base.RecordType)
	}

	return nil
}

func handleBounce(e emailx.PostmarkBounceEvent) {
	log.Printf("Bounce received for %s: %s (%s)", e.Email, e.Type, e.Description)
	// TODO: Update user status or subscription as needed
}

func handleSpamComplaint(e emailx.PostmarkSpamComplaintEvent) {
	log.Printf("Spam complaint received for %s", e.Email)
	// TODO: Unsubscribe user immediately
	// Logic: Find User by Email -> Set IsUnsubscribed = true (once field exists)
}
