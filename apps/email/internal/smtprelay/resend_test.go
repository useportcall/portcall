package smtprelay

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestSendViaResend(t *testing.T) {
	gotAuth := ""
	gotBody := ""
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotAuth = r.Header.Get("Authorization")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(strings.NewReader(`{"id":"ok"}`)),
			Header:     make(http.Header),
		}, nil
	})}
	s := &Session{
		backend: &Backend{cfg: Config{EmailProvider: "resend", ResendToken: "re_test", ResendAPI: "https://example.test/emails"}, httpClient: httpClient},
		from:    "sender@example.com",
		to:      []string{"recipient@example.com"},
		data:    []byte("Subject: Smoke Test\r\n\r\nhello"),
	}
	if err := s.sendViaResend(); err != nil {
		t.Fatalf("sendViaResend failed: %v", err)
	}
	if gotAuth != "Bearer re_test" {
		t.Fatalf("unexpected auth header: %q", gotAuth)
	}
	for _, want := range []string{`"subject":"Smoke Test"`, "recipient@example.com", "sender@example.com"} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("request body missing %q: %s", want, gotBody)
		}
	}
}
