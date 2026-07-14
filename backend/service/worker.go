package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
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

// Bounds how long a single message's processing can run — without this,
// a hung DB call or downstream dependency would permanently tie up one
// of the workerPoolSize goroutines draining the delivery channel.
const processDeliveryTimeout = 30 * time.Second

func (w *MessageWorker) processDelivery(d amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), processDeliveryTimeout)
	defer cancel()

	var env wire.Msg
	if err := json.Unmarshal(d.Body, &env); err != nil {
		slog.Error("worker: malformed message, discarding", "error", err)
		d.Nack(false, false) // don't requeue — malformed will never succeed
		return
	}

	processErr := w.Service.ProcessChat(ctx, &env)
	retryCount := retryCountFromHeaders(d.Headers)

	switch decideOutcome(processErr, retryCount) {
	case outcomeAck:
		d.Ack(false)

	case outcomeDiscard:
		slog.Error("worker: ProcessChat failed after max retries, discarding", "error", processErr, "retry_count", retryCount)
		d.Nack(false, false) // give up — stop requeuing this message

	case outcomeRetry:
		slog.Error("worker: ProcessChat failed, requeueing", "error", processErr, "retry_count", retryCount+1)
		// Re-publish with an incremented retry count (Nack's requeue can't carry
		// custom headers), then remove the original delivery from the queue —
		// but only once the replacement is safely enqueued.
		if pubErr := w.republishWithRetry(d.Body, amqp.Table{retryCountHeader: int32(retryCount + 1)}); pubErr != nil {
			slog.Error("worker: republish failed after retries, requeueing original instead", "error", pubErr)
			d.Nack(false, true)
			return
		}
		d.Ack(false)
	}
}

// republishRetryAttempts and republishRetryDelay bound how hard we try to get
// the incremented-retry copy onto the queue before giving up and requeuing the
// original untouched. This absorbs brief broker blips without ever needing the
// no-counter-increment fallback in processDelivery.
const republishRetryAttempts = 3
const republishRetryDelay = 100 * time.Millisecond

func (w *MessageWorker) republishWithRetry(body []byte, headers amqp.Table) error {
	var err error
	for i := 0; i < republishRetryAttempts; i++ {
		if err = w.MQ.PublishRaw(body, headers); err == nil {
			return nil
		}
		time.Sleep(republishRetryDelay)
	}
	return err
}

// deliveryOutcome is what a delivery should be resolved to. Kept separate
// from processDelivery's actual Ack/Nack/PublishRaw calls so the branching
// can be unit tested without a real AMQP channel or MessageService.
type deliveryOutcome int

const (
	outcomeAck     deliveryOutcome = iota // ProcessChat succeeded
	outcomeDiscard                        // failed and out of retries — drop it
	outcomeRetry                          // failed but retries remain — republish with count+1
)

// decideOutcome is the pure decision: given whether ProcessChat succeeded
// and how many times this delivery has already been retried, what happens
// to it next.
func decideOutcome(processErr error, retryCount int) deliveryOutcome {
	if processErr == nil {
		return outcomeAck
	}
	if retryCount >= maxDeliveryRetries {
		return outcomeDiscard
	}
	return outcomeRetry
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
