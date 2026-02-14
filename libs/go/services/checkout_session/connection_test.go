package checkout_session

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/gorm"
)

type connectionDB struct {
	dbx.IORM
	defaultConn  *models.Connection
	localConn    *models.Connection
	fallbackConn *models.Connection
}

func (m *connectionDB) FindForID(id uint, dest any) error {
	if m.defaultConn == nil {
		return gorm.ErrRecordNotFound
	}
	*(dest.(*models.Connection)) = *m.defaultConn
	return nil
}

func (m *connectionDB) FindFirst(dest any, conds ...any) error {
	query := conds[0].(string)
	if query == "app_id = ? AND source = ?" {
		if m.localConn == nil {
			return gorm.ErrRecordNotFound
		}
		*(dest.(*models.Connection)) = *m.localConn
		return nil
	}
	if m.fallbackConn == nil {
		return gorm.ErrRecordNotFound
	}
	*(dest.(*models.Connection)) = *m.fallbackConn
	return nil
}

func TestResolveConnection_DefaultConnectionWins(t *testing.T) {
	svc := &service{db: &connectionDB{defaultConn: &models.Connection{Source: "stripe"}}}
	conn, err := svc.resolveConnection(1, 10)
	if err != nil {
		t.Fatalf("resolveConnection() error = %v", err)
	}
	if conn.Source != "stripe" {
		t.Fatalf("expected stripe default connection, got %s", conn.Source)
	}
}

func TestResolveConnection_PrefersLocalFallback(t *testing.T) {
	svc := &service{db: &connectionDB{
		localConn:    &models.Connection{Source: "local"},
		fallbackConn: &models.Connection{Source: "stripe"},
	}}
	conn, err := svc.resolveConnection(1, 10)
	if err != nil {
		t.Fatalf("resolveConnection() error = %v", err)
	}
	if conn.Source != "local" {
		t.Fatalf("expected local fallback, got %s", conn.Source)
	}
}

func TestResolveConnection_NoConnectionConfigured(t *testing.T) {
	svc := &service{db: &connectionDB{}}
	_, err := svc.resolveConnection(1, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
