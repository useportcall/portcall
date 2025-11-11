package server

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
)

type IMultiplexer interface {
	HandleFunc(taskType string, handler HandlerFunc)
}

type multiplexer struct {
	instance *asynq.ServeMux
	db       dbx.IORM
	queue    qx.IQueue
	crypto   cryptox.ICrypto
}

func (m *multiplexer) HandleFunc(taskType string, handler HandlerFunc) {
	m.instance.HandleFunc(taskType, func(ctx context.Context, t *asynq.Task) error {
		c := &Context{
			Task:   t,
			orm:    m.db,     // Use the DB connection from the multiplexer
			queue:  m.queue,  // Use the Queue from the multiplexer
			crypto: m.crypto, // Use the Crypto from the multiplexer
		}

		return handler(c)
	})
}

type HandlerFunc = func(IContext) error

type IContext interface {
	DB() dbx.IORM
	Crypto() cryptox.ICrypto
	Queue() qx.IQueue
	Payload() []byte
}

type Context struct {
	Task   *asynq.Task
	orm    dbx.IORM
	queue  qx.IQueue
	crypto cryptox.ICrypto
}

func (c *Context) DB() dbx.IORM {
	return c.orm
}

func (c *Context) Queue() qx.IQueue {
	return c.queue
}

func (c *Context) Payload() []byte {
	return c.Task.Payload()
}

func (c *Context) Crypto() cryptox.ICrypto {
	return c.crypto
}
