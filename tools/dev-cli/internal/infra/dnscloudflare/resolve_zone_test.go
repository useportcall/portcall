package dnscloudflare

import (
	"net/http"
	"strings"
	"testing"
)

func TestResolveZoneIDFindsParentZone(t *testing.T) {
	client := testClient(func(req *http.Request) string {
		name := req.URL.Query().Get("name")
		if name == "portcall.com" {
			return `{"success":true,"errors":[],"result":[{"id":"zone-1","name":"portcall.com"}]}`
		}
		return `{"success":true,"errors":[],"result":[]}`
	})
	got, err := client.ResolveZoneID("api.dev.portcall.com", "")
	if err != nil {
		t.Fatalf("resolve zone id: %v", err)
	}
	if got != "zone-1" {
		t.Fatalf("unexpected zone id: %s", got)
	}
}

func TestResolveZoneIDUsesExplicitID(t *testing.T) {
	client := testClient(func(req *http.Request) string {
		t.Fatalf("unexpected request: %s", req.URL.String())
		return ""
	})
	got, err := client.ResolveZoneID("portcall.com", "zone-explicit")
	if err != nil {
		t.Fatalf("resolve zone id: %v", err)
	}
	if strings.TrimSpace(got) != "zone-explicit" {
		t.Fatalf("unexpected zone id: %s", got)
	}
}
