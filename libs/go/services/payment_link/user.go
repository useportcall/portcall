package payment_link

import (
	"net/mail"
	"strings"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *service) resolveUser(input *CreateInput) (*models.User, error) {
	if strings.TrimSpace(input.UserID) != "" {
		var user models.User
		if err := s.db.GetForPublicID(input.AppID, input.UserID, &user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	email := strings.ToLower(strings.TrimSpace(input.UserEmail))
	if email == "" {
		return nil, NewValidationError("one of user_id or user_email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, NewValidationError("invalid user_email")
	}
	var existing models.User
	err := s.db.FindFirst(&existing, "app_id = ? AND lower(email) = lower(?)", input.AppID, email)
	if err == nil {
		return &existing, nil
	}
	if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}

	user := &models.User{
		PublicID: dbx.GenPublicID("user"),
		AppID:    input.AppID,
		Email:    email,
		Name:     deriveUserName(input.UserName, email),
	}
	if err := s.db.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func deriveUserName(name, email string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed != "" {
		return trimmed
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) == 2 && strings.TrimSpace(parts[0]) != "" {
		return strings.TrimSpace(parts[0])
	}
	return "Customer"
}
