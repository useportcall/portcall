package followup

import "time"

const (
	defaultQuoteReminderAfter = 24 * time.Hour
	defaultQuoteFinalAfter    = 72 * time.Hour
)

const (
	StageQuoteWaiting  = "quote_waiting_for_payment"
	StageQuoteReminder = "quote_payment_reminder"
	StageQuoteFinal    = "quote_payment_final_notice"
)

type QuoteUnpaidInput struct {
	AcceptedAt       *time.Time
	Now              time.Time
	ReminderAfter    time.Duration
	FinalNoticeAfter time.Duration
}

type QuoteUnpaidDecision struct {
	ShouldNotify bool
	FinalNotice  bool
	Stage        string
}

func DecideQuoteUnpaid(input QuoteUnpaidInput) QuoteUnpaidDecision {
	if input.AcceptedAt == nil {
		return QuoteUnpaidDecision{Stage: StageQuoteWaiting}
	}
	now := input.Now
	if now.IsZero() || now.Before(*input.AcceptedAt) {
		return QuoteUnpaidDecision{Stage: StageQuoteWaiting}
	}
	reminderAfter, finalAfter := normalizeQuoteThresholds(input)
	age := now.Sub(*input.AcceptedAt)
	if age >= finalAfter {
		return QuoteUnpaidDecision{
			ShouldNotify: true,
			FinalNotice:  true,
			Stage:        StageQuoteFinal,
		}
	}
	if age >= reminderAfter {
		return QuoteUnpaidDecision{
			ShouldNotify: true,
			Stage:        StageQuoteReminder,
		}
	}
	return QuoteUnpaidDecision{Stage: StageQuoteWaiting}
}

func normalizeQuoteThresholds(input QuoteUnpaidInput) (time.Duration, time.Duration) {
	reminderAfter := input.ReminderAfter
	if reminderAfter <= 0 {
		reminderAfter = defaultQuoteReminderAfter
	}
	finalAfter := input.FinalNoticeAfter
	if finalAfter <= reminderAfter {
		finalAfter = defaultQuoteFinalAfter
	}
	return reminderAfter, finalAfter
}
