package service

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/leeryan2000/flashat/db"
	"github.com/leeryan2000/flashat/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

type routedDelivery struct {
	d   amqp.Delivery
	env wire.Msg
}

// used to determine which worker in the pool should handle a given conversation.
func workerIndexFor(conversationID string, poolSize int) int {
	h := fnv.New32a()
	h.Write([]byte(conversationID))
	return int(h.Sum32()) % poolSize
}

type MessageWorker struct {
	MQ      *db.RabbitMQClient
	Service *MessageService

	wg sync.WaitGroup
}

const workerPoolSize = 4

func NewMessageWorker(mq *db.RabbitMQClient, svc *MessageService) *MessageWorker {
	return &MessageWorker{
		MQ:      mq,
		Service: svc,
	}
}

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

	// Each worker gets its own buffered channel
	workerChans := make([]chan routedDelivery, workerPoolSize)
	for i := range workerChans {
		workerChans[i] = make(chan routedDelivery, 32)
	}

	slog.Info("worker: starting pool, waiting for messages", "pool_size", workerPoolSize)
	for i := range workerPoolSize {
		w.wg.Add(1)
		go func(ch <-chan routedDelivery) {
			defer w.wg.Done()
			for {
				select {
				case rd, ok := <-ch:
					if !ok {
						return
					}
					w.processDelivery(ctx, rd)
				case <-ctx.Done():
					slog.Info("worker: shutdown signal received, stopping")
					return
				}
			}
		}(workerChans[i])
	}

	// Single dispatcher: the only goroutine reading the shared RabbitMQ
	// delivery channel, so there's no race on read order.
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer func() {
			for _, ch := range workerChans {
				close(ch)
			}
		}()
		for {
			select {
			case d, ok := <-deliveries:
				if !ok {
					slog.Info("worker: delivery channel closed, shutting down")
					return
				}
				var env wire.Msg
				if err := json.Unmarshal(d.Body, &env); err != nil {
					slog.Error("worker: malformed message, discarding", "error", err)
					d.Nack(false, false)
					continue
				}
				idx := workerIndexFor(env.ConversationID, workerPoolSize)
				select {
				case workerChans[idx] <- routedDelivery{d: d, env: env}:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				slog.Info("worker: dispatcher shutdown signal received, stopping")
				return
			}
		}
	}()
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

func (w *MessageWorker) processDelivery(parent context.Context, rd routedDelivery) {
	ctx, cancel := context.WithTimeout(parent, processDeliveryTimeout)
	defer cancel()

	d, env := rd.d, rd.env

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
