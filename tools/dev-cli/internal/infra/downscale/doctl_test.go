package downscale

import "testing"

func TestParseDBRows(t *testing.T) {
	out := "id1 portcall-db pg online db-s-1vcpu-2gb 1\nid2 portcall-redis valkey online db-s-1vcpu-1gb 1"
	rows := ParseDBRows(out)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[1].Name != "portcall-redis" || rows[1].Engine != "valkey" || rows[1].Nodes != 1 {
		t.Fatalf("unexpected redis row: %+v", rows[1])
	}
}
