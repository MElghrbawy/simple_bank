package worker

import (
	"context"

	db "github.com/MElghrbawy/simple_bank/db/sqlc"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt *asynq.RedisClientOpt, store db.Store) TaskProcessor {

	server := asynq.NewServer(redisOpt, asynq.Config{})
	return &RedisTaskProcessor{server: server, store: store}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)

}
