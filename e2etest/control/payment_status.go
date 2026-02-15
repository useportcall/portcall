package control

import (
	"encoding/json"
	"net/http"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type setUserPaymentStatusRequest struct {
	UserID             string `json:"user_id"`
	SubscriptionStatus string `json:"subscription_status"`
	InvoiceStatus      string `json:"invoice_status"`
}

func (s *Server) handleSetUserPaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req setUserPaymentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := s.H.DB.FindFirst(&user, "app_id = ? AND public_id = ?", s.H.AppID, req.UserID); err != nil {
		http.Error(w, "user not found: "+err.Error(), http.StatusNotFound)
		return
	}
	if err := s.updateSubscriptionStatus(user.ID, req.SubscriptionStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.updateInvoiceStatus(user.ID, req.InvoiceStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (s *Server) updateSubscriptionStatus(userID uint, status string) error {
	if status == "" {
		return nil
	}
	subs := []models.Subscription{}
	if err := s.H.DB.ListWithOrderAndLimit(&subs, "created_at DESC", 1, "app_id = ? AND user_id = ?", s.H.AppID, userID); err != nil {
		return err
	}
	if len(subs) == 0 {
		return nil
	}
	subs[0].Status = status
	return s.H.DB.Save(&subs[0])
}

func (s *Server) updateInvoiceStatus(userID uint, status string) error {
	if status == "" {
		return nil
	}
	invoices := []models.Invoice{}
	if err := s.H.DB.ListWithOrderAndLimit(&invoices, "created_at DESC", 1, "app_id = ? AND user_id = ?", s.H.AppID, userID); err != nil {
		return err
	}
	if len(invoices) == 0 {
		return nil
	}
	invoices[0].Status = status
	return s.H.DB.Save(&invoices[0])
}
