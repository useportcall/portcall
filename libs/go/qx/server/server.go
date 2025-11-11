package server

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
)

type IServer interface {
	R()
	H(taskType string, handler HandlerFunc)
}

type server struct {
	instance *asynq.Server
	mux      *multiplexer
}

func (s *server) R() {
	if err := s.instance.Run(s.mux.instance); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s *server) H(taskType string, handler HandlerFunc) {
	s.mux.HandleFunc(taskType, handler)
}

func New(db dbx.IORM, crypto cryptox.ICrypto, queues map[string]int) IServer {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set")
	}

	q := qx.New()

	muxInstance := &multiplexer{asynq.NewServeMux(), db, q, crypto}

	instance := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Queues: queues,
		},
	)

	return &server{instance, muxInstance}
}
