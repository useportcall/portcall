package smtprelay

import (
	"io"
	"log"
	"os"
	"strings"
)

type smtpErrorFilterWriter struct {
	out io.Writer
}

func (w smtpErrorFilterWriter) Write(p []byte) (int, error) {
	msg := string(p)
	if isBenignSMTPNetworkError(msg) {
		return len(p), nil
	}
	return w.out.Write(p)
}

func isBenignSMTPNetworkError(msg string) bool {
	lower := strings.ToLower(msg)
	if !strings.Contains(lower, "error handling") {
		return false
	}
	patterns := []string{
		"connection reset by peer",
		"broken pipe",
		"use of closed network connection",
		": eof",
		"i/o timeout",
	}
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func newSMTPErrorLogger() *log.Logger {
	return log.New(smtpErrorFilterWriter{out: os.Stderr}, "smtp/server ", log.LstdFlags)
}
