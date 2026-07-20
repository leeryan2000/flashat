package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/leeryan2000/flashat/db"
	"github.com/leeryan2000/flashat/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageWorker owns all RabbitMQ consumer mechanics.
// It knows nothing about business logic — it just dequeues,
// calls the service, and translates the result into ack/nack.
type MessageWorker struct {
	MQ      *db.RabbitMQClient
	Service *MessageService

	wg sync.WaitGroup
}

// Number of goroutines concurrently draining the delivery channel.
// Bottleneck here is I/O (DB + websocket writes), not CPU, so this
// isn't tied to runtime.NumCPU().
const workerPoolSize = 4

func NewMessageWorker(mq *db.RabbitMQClient, svc *MessageService) *MessageWorker {
	return &MessageWorker{
		MQ:      mq,
		Service: svc,
	}
}

// Launch a pool of goroutines consuming from the same delivery channel.
// ctx is a shutdown signal, not a per-message context — cancelling it stops
// the pool from picking up new deliveries and bounds how long an in-flight
// one can keep running; see WaitForShutdown.
func (w *MessageWorker) Start(ctx context.Context) {
	deliveries, err := w.MQ.Ch.Consume(
		"chat_messages",
		"",    // consumer tag (auto-generated)
		false, // autoAck=false to allow explicit ack/nack
		false, false, false, nil,
	)
	if err != nil {
		slog.Error("worker: failed to start consumer", "error", err)
		os.Exit(1)
	}

	slog.Info("worker: starting pool, waiting for messages", "pool_size", workerPoolSize)
	for range workerPoolSize {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				select {
				case d, ok := <-deliveries:
					if !ok {
						slog.Info("worker: delivery channel closed, shutting down")
						return
					}
					w.processDelivery(ctx, d)
				case <-ctx.Done():
					slog.Info("worker: shutdown signal received, stopping")
					return
				}
			}
		}()
	}
}

// processDeliveryTimeout bounds a single ProcessChat attempt. Kept under the
// frontend's 10s ack-timeout (ChatContext.tsx) with a safety margin for
// network/broadcast latency
const processDeliveryTimeout = 8 * time.Second

const ShutdownTimeout = processDeliveryTimeout + 5*time.Second

// processRetryAttempts bounds how many times ProcessChat is retried in-process
// before the delivery is given up on.
const processRetryAttempts = 3
const processRetryDelay = 200 * time.Millisecond

func (w *MessageWorker) WaitForShutdown() {
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("worker: pool drained cleanly")
	case <-time.After(ShutdownTimeout):
		slog.Warn("worker: pool did not drain before timeout, continuing shutdown anyway")
	}
}

func (w *MessageWorker) processDelivery(parent context.Context, d amqp.Delivery) {
	ctx, cancel := context.WithTimeout(parent, processDeliveryTimeout)
	defer cancel()

	var env wire.Msg
	if err := json.Unmarshal(d.Body, &env); err != nil {
		slog.Error("worker: malformed message, discarding", "error", err)
		d.Nack(false, false) // don't requeue — malformed will never succeed
		return
	}

	var processErr error
retryLoop:
	for attempt := 1; attempt <= processRetryAttempts; attempt++ {
		processErr = w.Service.ProcessChat(ctx, &env)
		if processErr == nil {
			d.Ack(false)
			return
		}

		slog.Error("worker: ProcessChat failed", "error", processErr, "attempt", attempt)
		if attempt < processRetryAttempts {
			select {
			case <-time.After(processRetryDelay):
			case <-ctx.Done():
				break retryLoop // timeout/shutdown — stop retrying, fall through to discard
			}
		}
	}

	slog.Error("worker: ProcessChat failed after retries, discarding", "error", processErr)
	d.Nack(false, false) // no requeue — retries exhausted, this attempt is final
}
