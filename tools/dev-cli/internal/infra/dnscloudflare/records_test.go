package dnscloudflare

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestEnsureRecordsCheckOnlyReportsMissing(t *testing.T) {
	client := testClient(func(req *http.Request) string {
		if req.Method == "GET" {
			return `{"success":true,"errors":[],"result":[]}`
		}
		t.Fatalf("unexpected method: %s", req.Method)
		return ""
	})
	out, err := client.EnsureRecords("zone-1", []string{"api.portcall.com"}, "203.0.113.10", false)
	if err != nil {
		t.Fatalf("ensure records: %v", err)
	}
	if len(out) != 1 || out[0].Status != "missing" {
		t.Fatalf("unexpected outcome: %+v", out)
	}
}

func TestEnsureRecordsAutoCreatesAndUpdates(t *testing.T) {
	client := testClient(func(req *http.Request) string {
		switch {
		case req.Method == "GET" && strings.Contains(req.URL.RawQuery, "api.portcall.com"):
			return `{"success":true,"errors":[],"result":[]}`
		case req.Method == "GET" && strings.Contains(req.URL.RawQuery, "dashboard.portcall.com"):
			return `{"success":true,"errors":[],"result":[{"id":"rec-1","type":"A","name":"dashboard.portcall.com","content":"198.51.100.2"}]}`
		case req.Method == "POST":
			return `{"success":true,"errors":[],"result":{"id":"new-1"}}`
		case req.Method == "PUT":
			return `{"success":true,"errors":[],"result":{"id":"rec-1"}}`
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
			return ""
		}
	})
	out, err := client.EnsureRecords("zone-1", []string{"api.portcall.com", "dashboard.portcall.com"}, "203.0.113.10", true)
	if err != nil {
		t.Fatalf("ensure records: %v", err)
	}
	if out[0].Status != "created" || out[1].Status != "updated" {
		t.Fatalf("unexpected outcomes: %+v", out)
	}
}

func testClient(responder func(*http.Request) string) *Client {
	rt := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(responder(req))),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})
	hc := &http.Client{Transport: rt}
	return NewClientWithHTTP("token", "https://api.cloudflare.com/client/v4", hc)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
