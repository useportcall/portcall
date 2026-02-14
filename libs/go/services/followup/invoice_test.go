package followup_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/services/followup"
)

func TestDecideInvoiceFailure(t *testing.T) {
	cases := []struct {
		name    string
		input   followup.InvoiceFailureInput
		stage   string
		status  string
		final   bool
		attempt int
		max     int
	}{
		{
			name:    "first reminder",
			input:   followup.InvoiceFailureInput{Attempt: 1, MaxAttempts: 4},
			stage:   followup.StageInvoiceFirst,
			status:  "past_due",
			attempt: 1,
			max:     4,
		},
		{
			name:    "retry reminder",
			input:   followup.InvoiceFailureInput{Attempt: 2, MaxAttempts: 4},
			stage:   followup.StageInvoiceRetry,
			status:  "past_due",
			attempt: 2,
			max:     4,
		},
		{
			name:    "final attempt",
			input:   followup.InvoiceFailureInput{Attempt: 4, MaxAttempts: 4},
			stage:   followup.StageInvoiceFinal,
			status:  "uncollectible",
			final:   true,
			attempt: 4,
			max:     4,
		},
		{
			name:    "hard decline forces final",
			input:   followup.InvoiceFailureInput{Attempt: 1, MaxAttempts: 4, NoRetry: true},
			stage:   followup.StageInvoiceFinal,
			status:  "uncollectible",
			final:   true,
			attempt: 4,
			max:     4,
		},
		{
			name:    "invalid values normalize",
			input:   followup.InvoiceFailureInput{Attempt: 0, MaxAttempts: 0},
			stage:   followup.StageInvoiceFirst,
			status:  "past_due",
			attempt: 1,
			max:     4,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := followup.DecideInvoiceFailure(tc.input)
			if got.Stage != tc.stage || got.PaymentStatus != tc.status {
				t.Fatalf("unexpected decision: %+v", got)
			}
			if got.FinalAttempt != tc.final {
				t.Fatalf("expected final=%v got %v", tc.final, got.FinalAttempt)
			}
			if got.Attempt != tc.attempt || got.MaxAttempts != tc.max {
				t.Fatalf("expected attempt/max %d/%d got %d/%d", tc.attempt, tc.max, got.Attempt, got.MaxAttempts)
			}
		})
	}
}
