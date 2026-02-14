package feature

import (
"encoding/json"
"testing"
)

func TestCreateFeatureRequest_RequiresFeatureID(t *testing.T) {
	tests := []struct {
		name      string
		payload   string
		wantEmpty bool
	}{
		{
			name:      "correct field feature_id",
			payload:   `{"feature_id":"analytics","is_metered":false}`,
			wantEmpty: false,
		},
		{
			name:      "wrong field id instead of feature_id",
			payload:   `{"id":"analytics","is_metered":false}`,
			wantEmpty: true,
		},
		{
			name:      "empty body",
			payload:   `{}`,
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
var req CreateFeatureRequest
if err := json.Unmarshal([]byte(tt.payload), &req); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got := (req.FeatureID == ""); got != tt.wantEmpty {
				t.Errorf("FeatureID empty = %v, want %v", got, tt.wantEmpty)
			}
		})
	}
}
