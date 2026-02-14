package emailx

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestResendClientSend(t *testing.T) {
	gotAuth := ""
	gotBody := ""
	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			gotAuth = r.Header.Get("Authorization")
			b, _ := io.ReadAll(r.Body)
			gotBody = string(b)
			resp := &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader(`{"id":"ok"}`)),
				Header:     make(http.Header),
			}
			return resp, nil
		}),
	}
	client := &resendEmailClient{
		apiKey:     "re_test",
		apiURL:     "https://example.test/emails",
		httpClient: httpClient,
	}
	err := client.Send("<b>hello</b>", "test-subject", "hello@useportcall.com", []string{"test@example.com"})
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if gotAuth != "Bearer re_test" {
		t.Fatalf("expected auth header, got %q", gotAuth)
	}
	for _, want := range []string{"test-subject", "test@example.com", "hello@useportcall.com"} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("request missing %q: %s", want, gotBody)
		}
	}
}
