package main

import (
	"testing"
)

func TestSelectApps_AllAndNumbers(t *testing.T) {
	got, err := selectApps("all")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(deployApps) {
		t.Fatalf("expected %d apps, got %d", len(deployApps), len(got))
	}
	got, err = selectApps("1,3")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{deployApps[0].Name, deployApps[2].Name}
	if len(got) != len(want) {
		t.Fatalf("want %v got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("want %v got %v", want, got)
		}
	}
}

func TestSelectApps_Invalid(t *testing.T) {
	if _, err := selectApps("999"); err == nil {
		t.Fatal("expected error for invalid number")
	}
	if _, err := selectApps("unknown"); err == nil {
		t.Fatal("expected error for invalid app")
	}
}
