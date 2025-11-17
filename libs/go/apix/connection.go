package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Connection struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Connection) Set(connection *models.Connection) *Connection {
	c.ID = connection.PublicID
	c.Name = connection.Name
	c.Source = connection.Source
	c.PublicKey = connection.PublicKey
	c.CreatedAt = connection.CreatedAt
	c.UpdatedAt = connection.UpdatedAt
	return c
}
