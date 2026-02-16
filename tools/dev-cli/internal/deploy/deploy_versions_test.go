package deploy

import "testing"

func TestBumpVersion(t *testing.T) {
	cases := []struct {
		cur, bump, want string
	}{
		{"v1.2.3", "patch", "v1.2.4"},
		{"v1.2.3", "minor", "v1.3.0"},
		{"v1.2.3", "major", "v2.0.0"},
		{"v1.2.3", "skip", "v1.2.3"},
	}
	for _, tc := range cases {
		got, err := bumpVersion(tc.cur, tc.bump)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != tc.want {
			t.Fatalf("%s %s want %s got %s", tc.cur, tc.bump, tc.want, got)
		}
	}
}

func TestMappings(t *testing.T) {
	if deployName("billing") != "billing-worker" || helmValueName("billing") != "billingWorker" {
		t.Fatal("billing mapping mismatch")
	}
	if deployName("file") != "file-api" || helmValueName("file") != "fileApi" {
		t.Fatal("file mapping mismatch")
	}
	if deployName("dashboard") != "dashboard" || helmValueName("dashboard") != "dashboard" {
		t.Fatal("default mapping mismatch")
	}
	if deployName("email") != "email-worker" || helmValueName("email") != "emailWorker" {
		t.Fatal("email mapping mismatch")
	}
}
