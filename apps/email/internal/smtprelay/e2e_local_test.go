//go:build e2e_email_local

package smtprelay

import (
	"bytes"
	"io"
	"net"
	"net/http"
	netsmtp "net/smtp"
	"strings"
	"testing"
	"time"

	gsmtp "github.com/emersion/go-smtp"
)

type relayRoundTripFunc func(*http.Request) (*http.Response, error)

func (f relayRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestSMTPRelayLocalE2E(t *testing.T) {
	subjects := make(chan string, 1)
	client := &http.Client{Transport: relayRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), "Password Reset Local E2E") {
			subjects <- "Password Reset Local E2E"
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id":"ok"}`)),
			Header:     make(http.Header),
		}, nil
	})}
	be := &Backend{
		cfg:        Config{EmailProvider: "resend", ResendToken: "re_test", ResendAPI: "https://example.test"},
		httpClient: client,
	}
	s := gsmtp.NewServer(be)
	s.Domain = "localhost"
	s.AllowInsecureAuth = true
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() { _ = s.Serve(ln) }()

	msg := []byte("From: noreply@example.com\r\nTo: user@example.com\r\nSubject: Password Reset Local E2E\r\n\r\nUse code 1234\r\n")
	if err := netsmtp.SendMail(ln.Addr().String(), nil, "noreply@example.com", []string{"user@example.com"}, msg); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-subjects:
		if !strings.Contains(got, "Password Reset Local E2E") {
			t.Fatalf("unexpected subject: %s", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for resend request")
	}
}
