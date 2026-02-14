package server

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/logx"
	"github.com/useportcall/portcall/libs/go/qx"
)

type IServer interface {
	R() error
	H(taskType string, handler HandlerFunc)
	SetEmailClient(emailClient emailx.IEmailClient)
}

type server struct {
	instance *asynq.Server
	mux      *multiplexer
}

func (s *server) R() error { return s.instance.Run(s.mux.instance) }

func (s *server) H(taskType string, handler HandlerFunc) {
	s.mux.HandleFunc(taskType, handler)
}

func (s *server) SetEmailClient(emailClient emailx.IEmailClient) {
	s.mux.email = emailClient
}

func New(db dbx.IORM, crypto cryptox.ICrypto, queues map[string]int) (IServer, error) {
	logx.Init()

	q, err := qx.New()
	if err != nil {
		return nil, err
	}
	redisOpt, err := redisClientOptFromEnv()
	if err != nil {
		return nil, err
	}

	muxInstance := &multiplexer{asynq.NewServeMux(), db, q, crypto, nil}

	instance := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: queues,
		},
	)

	return &server{instance, muxInstance}, nil
}

func NewNoDeps(queues map[string]int) (IServer, error) {
	redisOpt, err := redisClientOptFromEnv()
	if err != nil {
		return nil, err
	}

	instance := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: queues,
		},
	)

	muxInstance := &multiplexer{asynq.NewServeMux(), nil, nil, nil, nil}

	return &server{instance, muxInstance}, nil
}

func redisClientOptFromEnv() (asynq.RedisClientOpt, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return asynq.RedisClientOpt{}, fmt.Errorf("REDIS_ADDR environment variable not set")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisTLS := os.Getenv("REDIS_TLS")

	if redisTLS == "true" {
		return asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}, nil
	}
	if redisPassword != "" {
		return asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
		}, nil
	}
	return asynq.RedisClientOpt{
		Addr: redisAddr,
	}, nil
}
