package server

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/logx"
	"github.com/useportcall/portcall/libs/go/qx"
)

type IServer interface {
	R()
	H(taskType string, handler HandlerFunc)
	SetEmailClient(emailClient emailx.IEmailClient)
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

func (s *server) SetEmailClient(emailClient emailx.IEmailClient) {
	s.mux.email = emailClient
}

func New(db dbx.IORM, crypto cryptox.ICrypto, queues map[string]int) IServer {
	logx.Init()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set")
	}

	q := qx.New()

	muxInstance := &multiplexer{asynq.NewServeMux(), db, q, crypto, nil}

	instance := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Queues: queues,
		},
	)

	return &server{instance, muxInstance}
}

func NewNoDeps(queues map[string]int) IServer {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set")
	}

	instance := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Queues: queues,
		},
	)

	muxInstance := &multiplexer{asynq.NewServeMux(), nil, nil, nil, nil}

	return &server{instance, muxInstance}
}
