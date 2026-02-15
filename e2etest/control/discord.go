package control

import (
	"encoding/json"
	"net/http"

	"github.com/useportcall/portcall/e2etest/harness"
)

func (s *Server) handleDiscordMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"mode": s.H.DiscordMode})
}

func (s *Server) handleDiscordMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.H.DiscordMode != harness.DiscordModeCapture {
		_ = json.NewEncoder(w).Encode(map[string]any{"mode": s.H.DiscordMode, "count": 0, "messages": []string{}})
		return
	}
	kind := r.URL.Query().Get("kind")
	cap := harness.CaptureForKind(s.H, kind)
	if cap == nil {
		http.Error(w, "kind must be signup or billing", http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"mode": s.H.DiscordMode, "count": cap.Count(), "messages": cap.Messages(),
	})
}

func (s *Server) handleDiscordReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.H.SignupWebhook != nil {
		s.H.SignupWebhook.Reset()
	}
	if s.H.BillingWebhook != nil {
		s.H.BillingWebhook.Reset()
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}
