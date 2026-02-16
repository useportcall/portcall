package dnscloudflare

import (
	"reflect"
	"testing"
)

func TestZoneCandidates(t *testing.T) {
	got := zoneCandidates("api.dev.portcall.com")
	want := []string{"api.dev.portcall.com", "dev.portcall.com", "portcall.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("zone candidates mismatch: got %v want %v", got, want)
	}
}
