package followup_test

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/services/followup"
)

func TestDecideQuoteUnpaid(t *testing.T) {
	now := time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC)
	accepted := now.Add(-36 * time.Hour)
	finalAccepted := now.Add(-96 * time.Hour)

	cases := []struct {
		name   string
		input  followup.QuoteUnpaidInput
		notify bool
		final  bool
		stage  string
	}{
		{
			name:  "missing accepted time",
			input: followup.QuoteUnpaidInput{Now: now},
			stage: followup.StageQuoteWaiting,
		},
		{
			name:  "before reminder window",
			input: followup.QuoteUnpaidInput{Now: now, AcceptedAt: ptrTime(now.Add(-4 * time.Hour))},
			stage: followup.StageQuoteWaiting,
		},
		{
			name:   "reminder window",
			input:  followup.QuoteUnpaidInput{Now: now, AcceptedAt: &accepted},
			notify: true,
			stage:  followup.StageQuoteReminder,
		},
		{
			name:   "final window",
			input:  followup.QuoteUnpaidInput{Now: now, AcceptedAt: &finalAccepted},
			notify: true,
			final:  true,
			stage:  followup.StageQuoteFinal,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := followup.DecideQuoteUnpaid(tc.input)
			if got.ShouldNotify != tc.notify || got.FinalNotice != tc.final || got.Stage != tc.stage {
				t.Fatalf("unexpected decision: %+v", got)
			}
		})
	}
}

func ptrTime(v time.Time) *time.Time { return &v }
