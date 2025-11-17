package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Secret struct {
	ID         string     `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
}

func (s *Secret) Set(secret *models.Secret) *Secret {
	s.ID = secret.PublicID
	s.CreatedAt = secret.CreatedAt
	s.UpdatedAt = secret.UpdatedAt
	s.DisabledAt = secret.DisabledAt
	return s
}
