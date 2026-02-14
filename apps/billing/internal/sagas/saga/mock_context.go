package saga

import (
	"context"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// mockContext implements server.IContext for integration tests.
// Embeds server.IContext so only the methods actually called by billing
// handlers need concrete implementations; all legacy methods on
// server.IContext that are not used will panic if called, which is
// the desired behaviour in tests.
type mockContext struct {
	server.IContext // embedded — provides default panicking stubs
	ctx            context.Context
	payload        []byte
	db             dbx.IORM
	queue          *RecordingQueue
	crypto         cryptox.ICrypto
}

func (c *mockContext) DB() dbx.IORM           { return c.db }
func (c *mockContext) Crypto() cryptox.ICrypto { return c.crypto }
func (c *mockContext) Queue() qx.IQueue       { return c.queue }
func (c *mockContext) Payload() []byte         { return c.payload }

// context.Context delegation — the embedded server.IContext provides the
// interface but we need concrete implementations for Deadline/Done/Err/Value.
func (c *mockContext) Done() <-chan struct{} { return c.ctx.Done() }
func (c *mockContext) Err() error           { return c.ctx.Err() }
func (c *mockContext) Value(key any) any    { return c.ctx.Value(key) }
