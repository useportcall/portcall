package control

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *Server) handleCheckoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	publicID := strings.TrimSpace(r.URL.Query().Get("public_id"))
	if publicID == "" {
		http.Error(w, "missing public_id", http.StatusBadRequest)
		return
	}
	var session models.CheckoutSession
	if err := s.H.DB.FindFirst(&session, "public_id = ?", publicID); err != nil {
		http.Error(w, "checkout session not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":                  session.PublicID,
		"external_session_id": session.ExternalSessionID,
	})
}
