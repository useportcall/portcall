package webhookx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx"
	"gorm.io/gorm"
)

func stripeSignatureHeader(secret string, payload []byte) string {
	ts := time.Now().Unix()
	signed := fmt.Sprintf("%d.%s", ts, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signed))
	return fmt.Sprintf("t=%d,v1=%s", ts, hex.EncodeToString(mac.Sum(nil)))
}

func connectionWithSecret(publicID, encrypted string) models.Connection {
	return models.Connection{PublicID: publicID, EncryptedWebhookSecret: &encrypted}
}

func braintreeConnection(publicID, encryptedKey, publicKey string) models.Connection {
	return models.Connection{
		PublicID:     publicID,
		Source:       "braintree",
		PublicKey:    publicKey,
		EncryptedKey: encryptedKey,
	}
}

type dbStub struct {
	dbx.IORM
	conn models.Connection
}

func (d *dbStub) FindFirst(dest any, conds ...any) error {
	c, ok := dest.(*models.Connection)
	if !ok {
		return gorm.ErrRecordNotFound
	}
	*c = d.conn
	return nil
}

type cryptoStub struct {
	decrypted string
}

func (c *cryptoStub) Encrypt(data string) (string, error)            { return data, nil }
func (c *cryptoStub) Decrypt(data string) (string, error)            { return c.decrypted, nil }
func (c *cryptoStub) CompareHash(hashed, plain string) (bool, error) { return true, nil }

type queueRecorder struct {
	qx.IQueue
	count int
}

func (q *queueRecorder) Enqueue(name string, payload any, queue string) error {
	q.count++
	return nil
}
