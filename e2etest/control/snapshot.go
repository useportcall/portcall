package control

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleSnapshotConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	url := s.H.SnapshotWebhookURL
	mode := "local"
	if url != "" {
		mode = "live"
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"mode":        mode,
		"webhook_url": url,
	})
}
