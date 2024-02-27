package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (d *RedisTaskDistributor) DistributeTaskSendEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(TaskSendVerifyEmail, payloadBytes, opts...)

	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).
		Msg("task enqueued")
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {

		return fmt.Errorf("could not unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := p.store.GetUser(ctx, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("could not get user: %w", err)
	}
	// send email to user

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processing task")
	return nil
}
