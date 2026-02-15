package smtprelay

import (
	"bytes"
	"testing"
)

func TestIsBenignSMTPNetworkError(t *testing.T) {
	tcs := []struct {
		name string
		msg  string
		want bool
	}{
		{"reset", "smtp/server error handling x: read: connection reset by peer", true},
		{"eof", "smtp/server error handling x: EOF", true},
		{"broken_pipe", "smtp/server error handling x: write: broken pipe", true},
		{"real_error", "smtp/server error handling x: invalid command sequence", false},
		{"other_log", "smtp/server accepted session", false},
	}
	for _, tc := range tcs {
		if got := isBenignSMTPNetworkError(tc.msg); got != tc.want {
			t.Fatalf("%s: got %v want %v", tc.name, got, tc.want)
		}
	}
}

func TestSMTPErrorFilterWriter(t *testing.T) {
	var out bytes.Buffer
	w := smtpErrorFilterWriter{out: &out}
	_, _ = w.Write([]byte("smtp/server error handling x: read: connection reset by peer\n"))
	if out.Len() != 0 {
		t.Fatal("expected benign network error to be suppressed")
	}
	_, _ = w.Write([]byte("smtp/server error handling x: invalid command sequence\n"))
	if out.Len() == 0 {
		t.Fatal("expected non-benign error to be written")
	}
}
