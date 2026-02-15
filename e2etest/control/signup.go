package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *Server) handlePrepareSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var app models.App
	if err := s.H.DB.FindFirst(&app, "id = ?", s.H.AppID); err != nil {
		http.Error(w, "seed app not found: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var account models.Account
	if err := s.H.DB.FindFirst(&account, "id = ?", app.AccountID); err != nil {
		http.Error(w, "seed account not found: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if account.Email == "admin@portcall.internal" {
		account.Email = fmt.Sprintf("seeded-%d@portcall.internal", time.Now().UnixNano())
		if err := s.H.DB.Save(&account); err != nil {
			http.Error(w, "save account failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "seed_account_email": account.Email})
}
