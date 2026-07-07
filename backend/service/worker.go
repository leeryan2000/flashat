package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

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
}

// Number of goroutines concurrently draining the delivery channel.
// Bottleneck here is I/O (DB + websocket writes), not CPU, so this
// isn't tied to runtime.NumCPU().
const workerPoolSize = 4

// Max times a failing message is retried before being dropped for good.
// Without this, a permanently-failing message (e.g. referencing a deleted
// conversation) requeues and fails forever, spinning the CPU and flooding logs.
const maxDeliveryRetries = 5

const retryCountHeader = "x-retry-count"

func NewMessageWorker(mq *db.RabbitMQClient, svc *MessageService) *MessageWorker {
	return &MessageWorker{
		MQ:      mq,
		Service: svc,
	}
}

// Launch a pool of goroutines consuming from the same delivery channel.
func (w *MessageWorker) Start() {
	deliveries, err := w.MQ.Ch.Consume(
		"chat_messages",
		"",    // consumer tag (auto-generated)
		false, // autoAck=false — we manually ACK after processing
		false, false, false, nil,
	)
	if err != nil {
		slog.Error("worker: failed to start consumer", "error", err)
		os.Exit(1)
	}

	slog.Info("worker: starting pool, waiting for messages", "pool_size", workerPoolSize)
	for range workerPoolSize {
		go func() {
			for d := range deliveries {
				w.processDelivery(d)
			}
			slog.Info("worker: delivery channel closed, shutting down")
		}()
	}
}

func (w *MessageWorker) processDelivery(d amqp.Delivery) {
	ctx := context.Background()

	var env wire.Msg
	if err := json.Unmarshal(d.Body, &env); err != nil {
		slog.Error("worker: malformed message, discarding", "error", err)
		d.Nack(false, false) // don't requeue — malformed will never succeed
		return
	}

	if err := w.Service.ProcessChat(ctx, &env); err != nil {
		retryCount := retryCountFromHeaders(d.Headers)
		if retryCount >= maxDeliveryRetries {
			slog.Error("worker: ProcessChat failed after max retries, discarding", "error", err, "retry_count", retryCount)
			d.Nack(false, false) // give up — stop requeuing this message
			return
		}

		slog.Error("worker: ProcessChat failed, requeueing", "error", err, "retry_count", retryCount+1)
		// Re-publish with an incremented retry count (Nack's requeue can't carry
		// custom headers), then remove the original delivery from the queue.
		if pubErr := w.MQ.PublishRaw(d.Body, amqp.Table{retryCountHeader: int32(retryCount + 1)}); pubErr != nil {
			slog.Error("worker: failed to republish for retry", "error", pubErr)
		}
		d.Ack(false)
		return
	}

	d.Ack(false)
}

func retryCountFromHeaders(headers amqp.Table) int {
	v, ok := headers[retryCountHeader]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case int32:
		return int(n)
	case int64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}
